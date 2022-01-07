package main

import (
	"flag"
	"fmt"
	"internal/update"
)

var Version = "development"
var Build = ""

func main() {
	updater := update.NewSelfUpdater(Version)

	versionFlag := false
	flag.BoolVar(&versionFlag, "version", false, "prints the version and exit")
	flag.Parse()

	if versionFlag {
		fmt.Println(updater.GetVersion())
		return
	}

	updater.StartUpdateCheck()

	// CMD HERE

	updater.Wait()
}
