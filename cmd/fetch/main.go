package fetch

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"strings"
)

var command = cli.Command{
	Name:    "fetch",
	Aliases: []string{"f"},
	Usage:   "Fetch a project from VCS",
	Action:  action,
}

func action(c *cli.Context) error {
	if _, ok := os.LookupEnv("SUDO_USER"); ok {
		panic("don't run fetch as sudo")
	}

	return runCommand()

	r, err := git.PlainClone("~/Projects/hub", false, &git.CloneOptions{
		URL:      "https://github.com/go-git/go-git",
		Progress: os.Stdout,
	})

	if err != nil {
		return err
	}

	return nil
}

func runCommand(dir string, c string, cmdline string, verbose bool, keyval ...string) ([]byte, error) {
	m := make(map[string]string)
	for i := 0; i < len(keyval); i += 2 {
		m[keyval[i]] = keyval[i+1]
	}

	args := strings.Fields(cmdline)
	for i, arg := range args {
		args[i] = expand(m, arg)
	}

	cmd := exec.Command(c, args...)
	cmd.Dir = dir
	//cmd.Env = base.AppendPWD(os.Environ(), cmd.Dir)
	cmd.Env = append(os.Environ(), "PWD="+cmd.Dir)

	//if cfg.BuildX {
	fmt.Fprintf(os.Stdout, "cd %s\n", dir)
	fmt.Fprintf(os.Stdout, "%s %s\n", c, strings.Join(args, " "))
	//}
	out, err := cmd.Output()
	if err != nil {
		//if verbose /*|| cfg.BuildV*/ {
		fmt.Fprintf(os.Stderr, "# cd %s; %s %s\n", dir, c, strings.Join(args, " "))
		if ee, ok := err.(*exec.ExitError); ok && len(ee.Stderr) > 0 {
			os.Stderr.Write(ee.Stderr)
		} else {
			fmt.Fprintf(os.Stderr, err.Error())
		}
		//}
	}

	return out, err
}

func Cmd() *cli.Command {
	return &command
}

// expand rewrites s to replace {k} with match[k] for each key k in match.
func expand(match map[string]string, s string) string {
	// We want to replace each match exactly once, and the result of expansion
	// must not depend on the iteration order through the map.
	// A strings.Replacer has exactly the properties we're looking for.
	oldNew := make([]string, 0, 2*len(match))
	for k, v := range match {
		oldNew = append(oldNew, "{"+k+"}", v)
	}
	return strings.NewReplacer(oldNew...).Replace(s)
}
