package main

import (
    "os"
    "text/tabwriter"

    "github-stats/lib/cli"
)

func main() {
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
    defer w.Flush()

    cli.Execute(w)
}
