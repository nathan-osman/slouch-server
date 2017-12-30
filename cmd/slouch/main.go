package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nathan-osman/slouch-server/server"
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

		// Initialize the server
		s, err := server.New(&server.Config{
			ClientID:     c.String("client-id"),
			ClientSecret: c.String("client-secret"),
			Addr:         c.String("server-addr"),
		})
		if err != nil {
			return err
		}
		defer s.Close()

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
