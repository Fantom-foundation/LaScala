package main

import (
	"os"
	"fmt"

	"github.com/urfave/cli/v2"
)

var RunRegistrarApp = cli.App {
	Name: "Amneris Registrar",
	Commands: []*cli.Command {
		&Serve,
		&Publish,
		&Subscribe,
	},
	Description: `
		Hello, world!
	`,
}

var Serve = cli.Command {
	Action: ServeRegistrar,
	Name: "serve",
	Usage: "",
	Flags: []cli.Flag {
		&ServePortFlag,
	},
}

var Publish = cli.Command {
	Action: PublishRegistrar,
	Name: "pub",
	Usage: "",
	Flags: []cli.Flag {
		&PublishAddressFlag,
		&PublishPortFlag,
		&PublishTopicFlag,
	},
}

var Subscribe = cli.Command {
	Action: SubscribeRegistrar,
	Name: "sub",
	Usage: "",
	Flags: []cli.Flag {
		&PublishAddressFlag,
		&PublishPortFlag,
		&PublishTopicFlag,
		&PublishMasterIdFlag,
	},
}


func main() {
	if err := RunRegistrarApp.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
