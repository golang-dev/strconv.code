package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:     "add",
			Aliases:  []string{"a"},
			Usage:    "add a task to the list",
			Category: "Add",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: ", c.Args().First())
				return nil
			},
		},
		{
			Name:     "template",
			Aliases:  []string{"t", "tmpl"},
			Usage:    "options for task templates",
			Category: "Template",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a new template",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing template",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}

	app.Name = "AppName"
	app.Usage = "application usage"
	app.Description = "application description"
	app.Version = "1.0.1"

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
