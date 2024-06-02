package main

import (
	"os"
	"fmt"

	"github.com/urfave/cli/v2"
)

var RunReporterApp = cli.App {
	Name: "Amneris Reporter",
	Commands: []*cli.Command {
		&Serve,
	},
	Description: `
		Hello, world!
	`,
}

var Serve = cli.Command {
	Action: ServeReporter,
	Name: "serve",
	Usage: "",
	Flags: []cli.Flag {
		&ServePortFlag,
		&StaticDirFlag,
		&ReportDbFlag,
		&RedisAddressFlag,
		&RedisPortFlag,
		&RedisTopicFlag,
	},
}

func main() {
	if err := RunReporterApp.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
