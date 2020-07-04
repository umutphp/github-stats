package app

import (
	ucli "github.com/urfave/cli/v2"
)

func New() *ucli.App {
	app := &ucli.App{
		Name:      "github-stats",
		Version:   "0.0.6",
		Usage:     "Get the total visit stats of your GitHub repositories",
		UsageText: "github-stats [global options]",
		Authors: []*ucli.Author{
			{Name: "Umut Işık", Email: "umutphp@gmail.com"},
		},
		Flags: []ucli.Flag{
			&ucli.IntFlag{
				Name:    "day",
				Aliases: []string{"d"},
				Usage:   "The number of days from today to show the stats",
				Value:   0,
			},
			&ucli.IntFlag{
				Name:    "show-details",
				Aliases: []string{"s"},
				Usage:   "Show detailed output or not. 0 to close. Default is 1",
				Value:   1,
			},
			&ucli.StringFlag{
				Name:     "token",
				Aliases:  []string{"t"},
				Value:    "",
				Usage:    "Personal access token got from GitHub to use the API",
				Required: true,
			},
		},
	}

	app.CustomAppHelpTemplate = `
NAME:
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}
USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}
VERSION:
   {{.Version}}{{end}}{{end}}{{if .Description}}
DESCRIPTION:
   {{.Description}}{{end}}{{if len .Authors}}
AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
   {{range $index, $author := .Authors}}{{if $index}}
   {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}
OPTIONS:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}
COPYRIGHT:
   {{.Copyright}}{{end}}

`

	return app
}
