package main

import (
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
	Log.SetLogLevel(Info)
	updater := update.NewSelfUpdater(Version)
	var versionFlag bool
	app := cli.NewApp()

	app.Usage = "Manage your development environment"
	app.Version = updater.GetVersion()
	app.EnableBashCompletion = true

	app.Commands = []*cli.Command{
		complete.Cmd(),
		//install.Cmd(),
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		Log.Fatal("app.run failed: ", err)
	}

	stop := false
	for _, i := range os.Args[1:] {
		if i == "--generate-bash-completion" {
			stop = true
			break
		} else if i == "--help" || i == "-help" || i == "help" {
			stop = true
			break
		} else if i == "-h" || i == "h" {
			stop = true
			break
		}
	}

	if versionFlag || stop {
		return
	}

	updater.StartUpdateCheck()

	updater.Wait()
}
