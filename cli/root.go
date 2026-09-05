package cli

import (
	"context"
	"github.com/urfave/cli/v3"
	"log"
	"os"
)

func Exec() {
	if err := peekCmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

var peekCmd = &cli.Command{
	Name:  "peek",
	Usage: "Just a simple peek app that gives me information about linux",
	Commands: []*cli.Command{
		cpuPeek,
		memPeek,
		stgPeek,
		portPeek,
	},
}
