package install

import (
	"github.com/urfave/cli/v2"
)

var command = cli.Command{
	Name:    "install",
	Aliases: []string{"i"},
	Usage:   "Install from a predefined list of suggested packages",
	Action:  action,
}

func action(c *cli.Context) error {

	return nil
}

func Cmd() *cli.Command {
	return &command
}
