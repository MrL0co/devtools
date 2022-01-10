package main

import (
	"fmt"
	"os"
	"sort"

	. "internal/logging"
	"internal/update"

	"github.com/urfave/cli/v2"
)

var Version = "development"
var Build = ""

func main() {
	Log.SetLogLevel(Debug)
	updater := update.NewSelfUpdater(Version)
	var versionFlag bool

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "version",
				Aliases:     []string{"v"},
				Usage:       "prints the version and exit",
				Destination: &versionFlag,
			},
		},
		Name:  "devtools",
		Usage: "Manage your development environment",
		Action: func(c *cli.Context) error {
			if c.Bool("version") {
				fmt.Println(updater.GetVersion())
				return nil
			}
			Log.Info("boom! I say!")

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "complete",
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
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		Log.Fatal("app.run failed: ", err)
	}

	if versionFlag {
		return
	}

	updater.StartUpdateCheck()

	updater.Wait()
}
