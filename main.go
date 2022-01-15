package main

import (
	"github.com/MrL0co/devtools/cmd/install"
	"os"
	"sort"

	"github.com/urfave/cli/v2"

	"github.com/MrL0co/devtools/cmd/complete"
	. "github.com/MrL0co/devtools/internal/logging"
	"github.com/MrL0co/devtools/internal/update"
)

var Version = "development"
var Build = ""

func main() {
	var updater *update.SelfUpdater

	app := cli.NewApp()

	app.Usage = "Manage your development environment"
	app.Version = Version
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:    "verbosity",
			Aliases: []string{"V"},
			Value:   3,
			Usage:   "set verbosity of logging to " + ListLogLevels(),
		},
	}

	app.Before = func(c *cli.Context) error {
		verbosity := LogLevel(c.Int("verbosity")).Valid()
		Log.SetLogLevel(verbosity)
		updater = update.NewSelfUpdater(Version)

		return nil
	}

	app.After = func(c *cli.Context) error {
		updater.StartUpdateCheck()

		updater.Wait()

		return nil
	}

	app.Commands = []*cli.Command{
		complete.Cmd(),
		install.Cmd(),
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		Log.Fatal("app.run failed: ", err)
	}
}
