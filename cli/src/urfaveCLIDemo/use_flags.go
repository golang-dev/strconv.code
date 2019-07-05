package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

var (
	flags []cli.Flag
	host  string
	port  int
)

func init() {
	flags = []cli.Flag{
		cli.StringFlag{
			Name:        "t, host, ip-address",
			Value:       "127.0.0.1",
			Usage:       "Server host",
			Destination: &host,
		},
		cli.IntFlag{
			Name:        "p, port",
			Value:       8000,
			Usage:       "Server port",
			Destination: &port,
		},
	}
}
func main() {
	app := cli.NewApp()
	app.Name = "AppName"
	app.Usage = "Application Usage"
	app.HideVersion = true
	app.Flags = flags
	app.Action = Action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Action(c *cli.Context) error {
	if c.Int("port") < 1024 {
		cli.ShowAppHelp(c)
		return cli.NewExitError("Ports below 1024 is not available", 2)
	}

	fmt.Printf("Listening at: http://%s:%d", host, c.Int("port"))
	return nil
}
