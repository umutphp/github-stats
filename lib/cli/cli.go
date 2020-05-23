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

type CLI struct {
	username string
	token string
	client *github.Client
	context context.Context
	accountTotalChannel chan int
    accountUniqueChannel chan int
    repoChannel chan *github.Repository
    haltChannel chan int
}

func New() CLI {
	cli := CLI{
		username: "",
		token: "07a4574f57fdbe125f37afab2264b64a9cde8d82",
	}

	cli.context  = context.Background()
    tokenSource := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: cli.token},
    )

    tokenClient := oauth2.NewClient(cli.context, tokenSource)

    cli.client               = github.NewClient(tokenClient)
    cli.accountTotalChannel  = make(chan int)
    cli.accountUniqueChannel = make(chan int)
    cli.repoChannel          = make(chan *github.Repository)
    cli.haltChannel          = make(chan int)

	return cli
}

func (cli *CLI) Initialize() bool {
	user, _, err := cli.client.Users.Get(cli.context, "")
    if err != nil {
        fmt.Printf("client.Users.Get() faled with '%s'\n", err)
        return false
    }
	
	cli.username = string(user.GetLogin())
	return true 
}

func (cli *CLI) GetRepos() []*github.Repository {
	// list all repositories for the authenticated user
    repos, _, err := cli.client.Repositories.List(cli.context, cli.username, nil)

    if err != nil {
        log.Fatal(err)
    }

    if len(repos) == 0 {
        log.Fatal(errors.New("User has no repositories"))
    }

    return repos
}

func (cli *CLI) RepoStat(w io.Writer) {
	for repo := range cli.repoChannel {
        fmt.Print(".")
		stats,_,_ := cli.client.Repositories.ListTrafficViews(cli.context, cli.username, repo.GetName(), nil)

        viewCount   := 0
        uniqueCount := 0
		for _,view := range stats.Views {
            if (view.GetTimestamp().After(time.Now().UTC().AddDate(0, 0, -1))) {
                uniqueCount += view.GetUniques()
                viewCount   += view.GetCount()
            }
		}

        fmt.Fprintf(w, "%s\t%d\t%d\t\n", repo.GetFullName(), viewCount, uniqueCount)
        cli.accountTotalChannel <- viewCount
        cli.accountUniqueChannel <- uniqueCount
	}

	cli.haltChannel <- 1
}

func (cli *CLI) Execute(w io.Writer) {
    repos     := cli.GetRepos()
    repoCount := len(repos)
    
    cli.accountTotalChannel  = make(chan int, repoCount)
    cli.accountUniqueChannel = make(chan int, repoCount)

	fmt.Fprintf(w, "Repository\tTotal View\tUnique View\t\n")

    fmt.Print("Checking repositories ")
    for i:=0;i<repoCount;i++ {
        go cli.RepoStat(w)
    }
    
	for _,repo := range repos {
		cli.repoChannel <- repo
	}
    
    close(cli.repoChannel)

    cli.Finiliaze(w, repoCount)
}

func (cli *CLI) Finiliaze(
	w io.Writer,
	repoCount int) {

	finishedCount := 0
    for finished := range cli.haltChannel {
    	finishedCount = finishedCount + finished

    	if finishedCount == repoCount {
    		close(cli.haltChannel)
            close(cli.accountTotalChannel)
            close(cli.accountUniqueChannel)
            fmt.Println("")
            fmt.Println("")
    	}
	}

    total := 0
    for t := range cli.accountTotalChannel {
        total += t
    }

    unique := 0
    for u := range cli.accountUniqueChannel {
        unique += u
    }

    fmt.Fprintf(w, "Total\t%d\t%d\t\n", total, unique)
}
