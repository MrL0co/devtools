package fetch

import (
	"errors"
	"fmt"
	"io/fs"
	urlpkg "net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type VcsCmd struct {
	Name string
	Cmd  string // name of binary to invoke command

	CreateCmd   []string // commands to download a fresh copy of a repository
	DownloadCmd []string // commands to download updates into an existing repository

	TagCmd         []tagCmd // commands to list tags
	TagLookupCmd   []tagCmd // commands to lookup tags before running tagSyncCmd
	TagSyncCmd     []string // commands to sync to specific tag
	TagSyncDefault []string // commands to sync to default tag

	Scheme  []string
	PingCmd string

	RemoteRepo  func(v *VcsCmd, rootDir string) (remoteRepo string, err error)
	ResolveRepo func(v *VcsCmd, rootDir, remoteRepo string) (realRepo string, err error)
}

// A tagCmd describes a command to list available tags
// that can be passed to tagSyncCmd.
type tagCmd struct {
	cmd     string // command to list tags
	pattern string // regexp to extract tags from list
}

// vcsGit describes how to use Git.
var vcsGit = &VcsCmd{
	Name: "Git",
	Cmd:  "git",

	CreateCmd:   []string{"clone -- {repo} {dir}"},
	DownloadCmd: []string{"pull --ff-only"},

	TagCmd: []tagCmd{
		// tags/xxx matches a git tag named xxx
		// origin/xxx matches a git branch named xxx on the default remote repository
		{"show-ref", `(?:tags|origin)/(\S+)$`},
	},
	TagLookupCmd: []tagCmd{
		{"show-ref tags/{tag} origin/{tag}", `((?:tags|origin)/\S+)$`},
	},
	TagSyncCmd: []string{"checkout {tag}"},
	// both createCmd and downloadCmd update the working dir.
	// No need to do more here. We used to 'checkout master'
	// but that doesn't work if the default branch is not named master.
	// DO NOT add 'checkout master' here.
	// See golang.org/issue/9032.
	TagSyncDefault: []string{"submodule update --init --recursive"},

	Scheme: []string{"git", "https", "http", "git+ssh", "ssh"},

	// Leave out the '--' separator in the ls-remote command: git 2.7.4 does not
	// support such a separator for that command, and this use should be safe
	// without it because the {scheme} value comes from the predefined list above.
	// See golang.org/issue/33836.
	PingCmd: "ls-remote {scheme}://{repo}",

	RemoteRepo: gitRemoteRepo,
}

// scpSyntaxRe matches the SCP-like addresses used by Git to access
// repositories by SSH.
var scpSyntaxRe = regexp.MustCompile(`^([a-zA-Z0-9_]+)@([a-zA-Z0-9._-]+):(.*)$`)

func gitRemoteRepo(vcsGit *VcsCmd, rootDir string) (remoteRepo string, err error) {
	cmd := "config remote.origin.url"
	errParse := errors.New("unable to parse output of git " + cmd)
	errRemoteOriginNotFound := errors.New("remote origin not found")
	outb, err := vcsGit.run1(rootDir, cmd, nil, false)
	if err != nil {
		// if it doesn't output any message, it means the config argument is correct,
		// but the config value itself doesn't exist
		if outb != nil && len(outb) == 0 {
			return "", errRemoteOriginNotFound
		}
		return "", err
	}
	out := strings.TrimSpace(string(outb))

	var repoURL *urlpkg.URL
	if m := scpSyntaxRe.FindStringSubmatch(out); m != nil {
		// Match SCP-like syntax and convert it to a URL.
		// Eg, "git@github.com:user/repo" becomes
		// "ssh://git@github.com/user/repo".
		repoURL = &urlpkg.URL{
			Scheme: "ssh",
			User:   urlpkg.User(m[1]),
			Host:   m[2],
			Path:   m[3],
		}
	} else {
		repoURL, err = urlpkg.Parse(out)
		if err != nil {
			return "", err
		}
	}

	// Iterate over insecure schemes too, because this function simply
	// reports the state of the repo. If we can't see insecure schemes then
	// we can't report the actual repo URL.
	for _, s := range vcsGit.Scheme {
		if repoURL.Scheme == s {
			return repoURL.String(), nil
		}
	}
	return "", errParse
}

// run runs the command line cmd in the given directory.
// keyval is a list of key, value pairs. run expands
// instances of {key} in cmd into value, but only after
// splitting cmd into individual arguments.
// If an error occurs, run prints the command line and the
// command's combined stdout+stderr to standard error.
// Otherwise run discards the command's output.
func (v *VcsCmd) run(dir string, cmd string, keyval ...string) error {
	_, err := v.run1(dir, cmd, keyval, true)
	return err
}

// runOutput is like run but returns the output of the command.
func (v *VcsCmd) runOutput(dir string, cmd string, keyval ...string) ([]byte, error) {
	return v.run1(dir, cmd, keyval, true)
}

// runVerboseOnly is like run but only generates error output to standard error in verbose mode.
func (v *VcsCmd) runVerboseOnly(dir string, cmd string, keyval ...string) error {
	_, err := v.run1(dir, cmd, keyval, false)
	return err
}

// run1 is the generalized implementation of run and runOutput.
func (v *VcsCmd) run1(dir string, cmdline string, keyval []string, verbose bool) ([]byte, error) {
	m := make(map[string]string)
	for i := 0; i < len(keyval); i += 2 {
		m[keyval[i]] = keyval[i+1]
	}
	args := strings.Fields(cmdline)
	for i, arg := range args {
		args[i] = expand(m, arg)
	}

	if len(args) >= 2 && args[0] == "-go-internal-mkdir" {
		var err error
		if filepath.IsAbs(args[1]) {
			err = os.Mkdir(args[1], fs.ModePerm)
		} else {
			err = os.Mkdir(filepath.Join(dir, args[1]), fs.ModePerm)
		}
		if err != nil {
			return nil, err
		}
		args = args[2:]
	}

	if len(args) >= 2 && args[0] == "-go-internal-cd" {
		if filepath.IsAbs(args[1]) {
			dir = args[1]
		} else {
			dir = filepath.Join(dir, args[1])
		}
		args = args[2:]
	}

	_, err := exec.LookPath(v.Cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"go: missing %s command. See https://golang.org/s/gogetcmd\n",
			v.Name)
		return nil, err
	}

	cmd := exec.Command(v.Cmd, args...)
	cmd.Dir = dir
	//cmd.Env = base.AppendPWD(os.Environ(), cmd.Dir)
	cmd.Env = append(os.Environ(), "PWD="+cmd.Dir)

	//if cfg.BuildX {
	//	fmt.Fprintf(os.Stderr, "cd %s\n", dir)
	//	fmt.Fprintf(os.Stderr, "%s %s\n", v.Cmd, strings.Join(args, " "))
	//}
	out, err := cmd.Output()
	if err != nil {
		if verbose /*|| cfg.BuildV*/ {
			fmt.Fprintf(os.Stderr, "# cd %s; %s %s\n", dir, v.Cmd, strings.Join(args, " "))
			if ee, ok := err.(*exec.ExitError); ok && len(ee.Stderr) > 0 {
				os.Stderr.Write(ee.Stderr)
			} else {
				fmt.Fprintf(os.Stderr, err.Error())
			}
		}
	}
	return out, err
}

//
//// expand rewrites s to replace {k} with match[k] for each key k in match.
//func expand(match map[string]string, s string) string {
//	// We want to replace each match exactly once, and the result of expansion
//	// must not depend on the iteration order through the map.
//	// A strings.Replacer has exactly the properties we're looking for.
//	oldNew := make([]string, 0, 2*len(match))
//	for k, v := range match {
//		oldNew = append(oldNew, "{"+k+"}", v)
//	}
//	return strings.NewReplacer(oldNew...).Replace(s)
//}

// Ping pings to determine scheme to use.
func (v *VcsCmd) Ping(scheme, repo string) error {
	return v.runVerboseOnly(".", v.PingCmd, "scheme", scheme, "repo", repo)
}

// Create creates a new copy of repo in dir.
// The parent of dir must exist; dir must not.
func (v *VcsCmd) Create(dir, repo string) error {
	for _, cmd := range v.CreateCmd {
		if err := v.run(".", cmd, "dir", dir, "repo", repo); err != nil {
			return err
		}
	}
	return nil
}

// Download downloads any new changes for the repo in dir.
func (v *VcsCmd) Download(dir string) error {
	for _, cmd := range v.DownloadCmd {
		if err := v.run(dir, cmd); err != nil {
			return err
		}
	}
	return nil
}

// Tags returns the list of available tags for the repo in dir.
func (v *VcsCmd) Tags(dir string) ([]string, error) {
	var tags []string
	for _, tc := range v.TagCmd {
		out, err := v.runOutput(dir, tc.cmd)
		if err != nil {
			return nil, err
		}
		re := regexp.MustCompile(`(?m-s)` + tc.pattern)
		for _, m := range re.FindAllStringSubmatch(string(out), -1) {
			tags = append(tags, m[1])
		}
	}
	return tags, nil
}

// tagSync syncs the repo in dir to the named tag,
// which either is a tag returned by tags or is v.tagDefault.
func (v *VcsCmd) TagSync(dir, tag string) error {
	if v.TagSyncCmd == nil {
		return nil
	}
	if tag != "" {
		for _, tc := range v.TagLookupCmd {
			out, err := v.runOutput(dir, tc.cmd, "tag", tag)
			if err != nil {
				return err
			}
			re := regexp.MustCompile(`(?m-s)` + tc.pattern)
			m := re.FindStringSubmatch(string(out))
			if len(m) > 1 {
				tag = m[1]
				break
			}
		}
	}

	if tag == "" && v.TagSyncDefault != nil {
		for _, cmd := range v.TagSyncDefault {
			if err := v.run(dir, cmd); err != nil {
				return err
			}
		}
		return nil
	}

	for _, cmd := range v.TagSyncCmd {
		if err := v.run(dir, cmd, "tag", tag); err != nil {
			return err
		}
	}
	return nil
}
