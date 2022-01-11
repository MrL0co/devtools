package main

import (
	"flag"
	"fmt"
	"github.com/posener/complete"
	"os"
	"sort"
	"strings"

	. "internal/logging"
	"internal/update"

	"github.com/urfave/cli/v2"
)

var Version = "development"
var Build = ""

func PrintCompletions(app *cli.App) {
	var cmpl []string
	for _, appFlag := range app.Flags {
		if documentFlag, ok := appFlag.(cli.DocGenerationFlag); ok {
			for _, name := range documentFlag.Names() {
				cmpl = append(
					cmpl, fmt.Sprintf(" '-%s[%s]'", name, documentFlag.GetUsage()))
			}
		}

	}
	fmt.Print(
		fmt.Sprintf("_arguments -s %s", strings.Join(cmpl, " ")),
	)
}

func main() {
	Log.SetLogLevel(Info)
	updater := update.NewSelfUpdater(Version)
	var versionFlag bool
	app := cli.NewApp()
	app.Usage = "Manage your development environment"

	//app.EnableBashCompletion = true

	// create the complete command
	cmp := complete.New(
		app.Name,
		complete.Command{Flags: complete.Flags{"-name": complete.PredictAnything}},
	)

	// AddFlags adds the completion flags to the program flags,
	// in case of using non-default flag set, it is possible to pass
	// it as an argument.
	// it is possible to set custom flags name
	// so when one will type 'self -h', he will see '-complete' to install the
	// completion and -uncomplete to uninstall it.
	cmp.CLI.InstallName = "complete"
	cmp.CLI.UninstallName = "uncomplete"
	cmp.AddFlags(nil)

	// parse the flags - both the program's flags and the completion flags
	flag.Parse()

	// run the completion, in case that the completion was invoked
	// and ran as a completion script or handled a flag that passed
	// as argument, the Run method will return true,
	// in that case, our program have nothing to do and should return.
	if cmp.Complete() {
		return
	}

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "version",
			Aliases:     []string{"v"},
			Usage:       "prints the version and exit",
			Destination: &versionFlag,
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.Bool("version") {
			fmt.Println(updater.GetVersion())
			return nil
		}
		Log.Info("boom! I say!")

		return nil
	}

	stop := false
	app.Commands = []*cli.Command{
		{
			Name:  "complete",
			Usage: "use for zsh auto completion",
			Action: func(c *cli.Context) error {
				PrintCompletions(app)
				stop = true
				return nil
			},
		},
		{
			Name:    "init",
			Aliases: []string{"c"},
			Usage:   "complete a task on the list",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a task to the list",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		Log.Fatal("app.run failed: ", err)
	}

	help := flag.Bool("help", false, "")
	h := flag.Bool("h", false, "")

	if versionFlag || *help || *h || stop {
		return
	}

	updater.StartUpdateCheck()

	updater.Wait()
}
