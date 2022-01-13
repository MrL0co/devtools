package main

import (
	"devtools/cmd/complete"
	"devtools/cmd/install"
	"flag"
	"fmt"
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

	app.EnableBashCompletion = true

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

	help := flag.Bool("help", false, "")
	h := flag.Bool("h", false, "")

	if versionFlag || *help || *h {
		return
	}
	//
	//updater.StartUpdateCheck()
	//
	//updater.Wait()
}
