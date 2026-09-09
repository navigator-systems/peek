package cli

import (
	"context"

	"peek/cmds"

	"github.com/urfave/cli/v3"
)

var cpuPeek = &cli.Command{
	Name:  "cpu",
	Usage: "peek cpu",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		cmds.CPU()
		return nil
	},
}

var memPeek = &cli.Command{
	Name:  "mem",
	Usage: "peek mem",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		cmds.MEM()
		return nil
	},
}

var stgPeek = &cli.Command{
	Name:  "stg",
	Usage: "peek stg",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		cmds.STG()
		return nil
	},
}

var portPeek = &cli.Command{
	Name:  "port",
	Usage: "peek port",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		portFlag := cmd.Args().Get(0)
		cmds.PORT(portFlag)
		return nil
	},
}

var procPeek = &cli.Command{
	Name:  "proc",
	Usage: "peek proc",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		actionFlag := cmd.Args().Get(0)
		cmds.PROC(actionFlag)
		return nil
	},
}
