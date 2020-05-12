package main

import (
	"fmt"
	"context"

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

    for i:=0;i<MAX_GOR_COUNT;i++ {
    	go repoStat(repoChannel, halt, client, ctx)
    }
    
    if err == nil {
    	for _,repo := range repos {
    		repoChannel <- repo
    	}
    }

    close(repoChannel)

    finishedCount := 0
    for finished := range halt {
    	finishedCount = finishedCount + finished

    	if finishedCount == MAX_GOR_COUNT {
    		close(halt)
    	}
	}
}

func repoStat(c chan *github.Repository, halt chan int, client *github.Client, ctx context.Context) {
	for repo := range c {
		stats,_,_ := client.Repositories.ListTrafficViews(ctx, "umutphp", repo.GetName(), nil)
		fmt.Printf("%s\n", repo.GetFullName())
		fmt.Println("")
		fmt.Println(stats.GetCount(), stats.GetUniques())

		for _,view := range stats.Views {
			fmt.Println(view.GetTimestamp(), view.GetCount(), view.GetUniques())
		}
	}

	halt <- 1
}
