package main

import (
	"fmt"
	"context"
    "time"
    "log"
    "errors"

	"golang.org/x/oauth2"
	"github.com/google/go-github/v31/github"	// with go modules enabled (GO111MODULE=on or outside GOPATH)
)


const MAX_GOR_COUNT = 4

func main() {
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

    accountTotal  := make(chan int, len(repos))
    accountUnique := make(chan int, len(repos))

    for i:=0;i<MAX_GOR_COUNT;i++ {
    	go repoStat(accountTotal, accountUnique, repoChannel, halt, client, ctx)
    }
    
	for _,repo := range repos {
		repoChannel <- repo
	}
    
    close(repoChannel)

    finishedCount := 0
    for finished := range halt {
    	finishedCount = finishedCount + finished

    	if finishedCount == MAX_GOR_COUNT {
    		close(halt)
            close(accountTotal)
            close(accountUnique)
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

    fmt.Println("Total", total, unique)
}

func repoStat(
    accountTotal chan int,
    accountUnique chan int,
    c chan *github.Repository,
    halt chan int,
    client *github.Client,
    ctx context.Context) {

	for repo := range c {
		stats,_,_ := client.Repositories.ListTrafficViews(ctx, "umutphp", repo.GetName(), nil)

        viewCount   := 0
        uniqueCount := 0
		for _,view := range stats.Views {
            if (view.GetTimestamp().After(time.Now().AddDate(0, 0, -1))) {
                uniqueCount += view.GetUniques()
                viewCount   += view.GetCount()
            }
		}

        fmt.Println(repo.GetFullName(), viewCount, uniqueCount)
        accountTotal <- viewCount
        accountUnique <- uniqueCount
	}

	halt <- 1
}
