package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "slouch"
	app.Usage = "Server for Slouch that facilitates login"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "client-id",
			EnvVar: "CLIENT_ID",
			Usage:  "Slack client ID",
		},
		cli.StringFlag{
			Name:   "client-secret",
			EnvVar: "CLIENT_SECRET",
			Usage:  "Slack client secret",
		},
		cli.StringFlag{
			Name:   "server-addr",
			Value:  ":8000",
			EnvVar: "SERVER_ADDR",
			Usage:  "server address",
		},
	}
	app.Action = func(c *cli.Context) error {

		// Wait for SIGINT or SIGTERM
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		return nil
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err.Error())
	}
}
