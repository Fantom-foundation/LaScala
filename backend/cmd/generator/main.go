package main

import (
	"os"
	"fmt"

	"github.com/urfave/cli/v2"
)

var RunWorkerApp = cli.App {
	Name: "Queue Worker",
	Commands: []*cli.Command {
		&Peek,
		&Push,
		&Pop,
	},
	Description: `
		Hello, world!
	`,
}

var Peek = cli.Command {
	Action: WorkerPeek,
	Name: "peek",
	Usage: "",
	Flags: []cli.Flag {
		&WorkerAddressFlag,
		&WorkerPortFlag,
		&WorkerTopicFlag,
	},
}

var Push = cli.Command {
	Action: WorkerPush,
	Name: "push",
	Usage: "",
	Flags: []cli.Flag {
		&WorkerAddressFlag,
		&WorkerPortFlag,
		&WorkerTopicFlag,
		&WorkerTypFlag,
		&WorkerMasterFlag,
		&WorkerRunFlag,
	},
}

var Pop = cli.Command {
	Action: WorkerPop,
	Name: "pop",
	Usage: "",
	Flags: []cli.Flag {
		&WorkerAddressFlag,
		&WorkerPortFlag,
		&WorkerTopicFlag,
	},
}

func main() {
	if err := RunWorkerApp.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
