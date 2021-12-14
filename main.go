package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

const (
	NAME    = "pango"
	USAGE   = "Paranoid text spacing"
	VERSION = "0.8.0"
)

var (
	errChan = make(chan error)
)

func main() {
	app := cli.NewApp()
	app.Name = NAME
	app.Usage = USAGE
	app.Version = VERSION
	app.Commands = []cli.Command{
		{
			Name:    "text",
			Usage:   "Performs paranoid text spacing on text",
			Aliases: []string{"t"},
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					_ = cli.ShowSubcommandHelp(c)
					return
				}

				text := c.Args().First()
				fmt.Println(SpacingText(text))
			},
		},
		{
			Name:    "file",
			Usage:   "Performs paranoid text spacing on files",
			Aliases: []string{"f"},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "write, w",
					Usage: fmt.Sprintf(`write result to (source) file instead of stdout of specific output file name`),
				},
				cli.BoolFlag{
					Name:  "comments, c",
					Usage: fmt.Sprintf(`only modify the codes' comments or not`),
				},
			},
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					_ = cli.ShowSubcommandHelp(c)
					return
				}

				writeOnly := c.Bool("write")
				commentsOnly := c.Bool("comments")

				if len(c.Args()) > 2 && writeOnly {
					color.Red(`can't use the "-write" flag with multiple files`)
					os.Exit(1)
				}

				go func() {
					for range c.Args() {
						err := <-errChan
						if err != nil {
							color.Red("%s", err)
						}
					}
				}()

				if c.Args()[0] == "." || strings.HasSuffix(c.Args()[0], "/") {
					processDir(c.Args()[0], writeOnly, commentsOnly)
				} else {
					for _, filename := range c.Args() {
						processFile(filename, writeOnly, commentsOnly)
					}
				}
			},
		},
	}

	_ = app.Run(os.Args)
}
