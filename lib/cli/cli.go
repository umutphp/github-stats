package cli

import (
	"fmt"
	"context"
	"log"
	"errors"
	"io"
	"time"

	"golang.org/x/oauth2"
	"github.com/google/go-github/v31/github"
)

func GetRepos(username string) []*github.Repository {
	ctx := context.Background()
    ts  := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: "07a4574f57fdbe125f37afab2264b64a9cde8d82"},
    )

    tc          := oauth2.NewClient(ctx, ts)
    client      := github.NewClient(tc)

	// list all repositories for the authenticated user
    repos, _, err := client.Repositories.List(ctx, username, nil)

    if err != nil {
        log.Fatal(err)
    }

    if len(repos) == 0 {
        log.Fatal(errors.New("User has no repositories"))
    }

    return repos
}

func RepoStat(
    w io.Writer,
    accountTotalChannel chan int,
    accountUniqueChannel chan int,
    repoChannel chan *github.Repository,
    haltChannel chan int) {

    ctx := context.Background()
    ts  := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: "07a4574f57fdbe125f37afab2264b64a9cde8d82"},
    )

    tc          := oauth2.NewClient(ctx, ts)
    client      := github.NewClient(tc)

	for repo := range repoChannel {
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
        accountTotalChannel <- viewCount
        accountUniqueChannel <- uniqueCount
	}

	haltChannel <- 1
}

func Execute(w io.Writer) {
    repos     := GetRepos("umutphp")
    repoCount := len(repos)
    accountTotalChannel  := make(chan int, repoCount)
    accountUniqueChannel := make(chan int, repoCount)
    repoChannel := make(chan *github.Repository)
    haltChannel := make(chan int)

	fmt.Fprintf(w, "Repository\tTotal View\tUnique View\t\n")

    fmt.Print("Checking repositories ")
    for i:=0;i<repoCount;i++ {
        go RepoStat(w, accountTotalChannel, accountUniqueChannel, repoChannel, haltChannel)
    }
    
	for _,repo := range repos {
		repoChannel <- repo
	}
    
    close(repoChannel)

    Finiliaze(w, repoCount, accountTotalChannel, accountUniqueChannel, haltChannel)
}

func Finiliaze(
	w io.Writer,
	repoCount int,
	accountTotalChannel chan int,
	accountUniqueChannel chan int,
	haltChannel chan int) {

	finishedCount := 0
    for finished := range haltChannel {
    	finishedCount = finishedCount + finished

    	if finishedCount == repoCount {
    		close(haltChannel)
            close(accountTotalChannel)
            close(accountUniqueChannel)
            fmt.Println("")
            fmt.Println("")
    	}
	}

    total := 0
    for t := range accountTotalChannel {
        total += t
    }

    unique := 0
    for u := range accountUniqueChannel {
        unique += u
    }

    fmt.Fprintf(w, "Total\t%d\t%d\t\n", total, unique)
}
