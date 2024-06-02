package main

import (
	"fmt"
	"net/http"
	"io"

	"github.com/urfave/cli/v2"
)

var (
	ServePortFlag = cli.IntFlag {
		Name: 	"port",
		Usage: 	"port number to use",
		Aliases: []string{"p"},
		Value: 3333,
	}
)

func ServeRegistrar(ctx *cli.Context) error {
	return serveRegistrar(
		ctx.Int("port"),
	)
}

func serveRegistrar(port int) error {
	http.HandleFunc("/register", registerRun)
	
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func registerRun(w http.ResponseWriter, r *http.Request) {
	const (
		addr = "127.0.0.1"
		port = 6379
		topic = "aida"
	)
	io.WriteString(w, fmt.Sprintf("Starting a consumer at %s:%d < %s!\n", addr, port, topic))
}
