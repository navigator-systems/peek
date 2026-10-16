package cli

import (
	"github.com/urfave/cli/v3"
)

var jsonFlag = &cli.BoolFlag{
	Name:  "json",
	Value: false,
	Usage: "output in json format",
}
