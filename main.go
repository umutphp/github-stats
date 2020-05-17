package main

import (
	"fmt"
	"context"
    "time"
    "log"
    "errors"
    "os"
    "text/tabwriter"
    "io"

	"golang.org/x/oauth2"
	"github.com/google/go-github/v31/github"
)

func main() {
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
    defer w.Flush()
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: "07a4574f57fdbe125f37afab2264b64a9cde8d82"},
    )

    tc          := oauth2.NewClient(ctx, ts)
    repoChannel := make(chan *github.Repository)
    client      := github.NewClient(tc)
    halt        := make(chan int)

    // list all repositories for the authenticated user
    repos, _, err := client.Repositories.List(ctx, "umutphp", nil)

    if err != nil {
        log.Fatal(err)
    }

    if len(repos) == 0 {
        log.Fatal(errors.New("umutphp has no repositories"))
    }

    repoCount     := len(repos)
    accountTotal  := make(chan int, repoCount)
    accountUnique := make(chan int, repoCount)

    fmt.Fprintf(w, "Repository\tTotal View\tUnique View\t\n")

    fmt.Print("Checking repositories ")
    for i:=0;i<repoCount;i++ {
    	go repoStat(w, accountTotal, accountUnique, repoChannel, halt, client, ctx)
    }
    
	for _,repo := range repos {
		repoChannel <- repo
	}
    
    close(repoChannel)

    finishedCount := 0
    for finished := range halt {
    	finishedCount = finishedCount + finished

    	if finishedCount == repoCount {
    		close(halt)
            close(accountTotal)
            close(accountUnique)
            fmt.Println("")
            fmt.Println("")
    	}
	}

    total := 0
    for t := range accountTotal {
        total += t
    }

    unique := 0
    for u := range accountUnique {
        unique += u
    }

    fmt.Fprintf(w, "Total\t%d\t%d\t\n", total, unique)
}

func repoStat(
    w io.Writer,
    accountTotal chan int,
    accountUnique chan int,
    c chan *github.Repository,
    halt chan int,
    client *github.Client,
    ctx context.Context) {

	for repo := range c {
        fmt.Print(".")
		stats,_,_ := client.Repositories.ListTrafficViews(ctx, "umutphp", repo.GetName(), nil)

        viewCount   := 0
        uniqueCount := 0
		for _,view := range stats.Views {
            if (view.GetTimestamp().After(time.Now().UTC().AddDate(0, 0, -1))) {
                uniqueCount += view.GetUniques()
                viewCount   += view.GetCount()
            }
		}

        fmt.Fprintf(w, "%s\t%d\t%d\t\n", repo.GetFullName(), viewCount, uniqueCount)
        accountTotal <- viewCount
        accountUnique <- uniqueCount
	}

	halt <- 1
}
