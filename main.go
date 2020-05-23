package main

import (
    "os"
    "text/tabwriter"

    "github-stats/lib/cli"
)

func main() {
    CLI := cli.New()
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
    defer w.Flush()

    if CLI.Initialize() == true {
        CLI.Execute(w)
    }
}
