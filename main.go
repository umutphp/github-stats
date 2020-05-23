package main

import (
    "os"
    "text/tabwriter"
    "fmt"

    "github-stats/lib/cli"
    "github-stats/lib/app"

    ucli "github.com/urfave/cli/v2"
)

func main() {
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
    defer w.Flush()
    application := app.New()

    application.Action = func(c *ucli.Context) error {
        CLI := cli.New(c.String("token"))

        CLI.SetDay(c.Int("day"))

        if c.Int("verbose") == 0 {
            CLI.SetVerbose(c.Int("verbose"))
        }

        if CLI.Initialize() == true {
            CLI.Execute(w)
        }
        return nil
    }

    err := application.Run(os.Args)

    if err != nil {
        fmt.Println(err)
        return
    }
}
