package cli

import (
	"context"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
)

type CLI struct {
	username             string
	token                string
	day                  int
	showDetails          int
	client               *github.Client
	context              context.Context
	accountTotalChannel  chan int
	accountUniqueChannel chan int
	repoChannel          chan *github.Repository
	haltChannel          chan int
}

func New(token string) CLI {
	cli := CLI{
		username:    "",
		token:       "",
		day:         0,
		showDetails: 1,
	}

	cli.SetToken(token)

	cli.context = context.Background()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cli.token},
	)

	tokenClient := oauth2.NewClient(cli.context, tokenSource)

	cli.client = github.NewClient(tokenClient)
	cli.accountTotalChannel = make(chan int)
	cli.accountUniqueChannel = make(chan int)
	cli.repoChannel = make(chan *github.Repository)
	cli.haltChannel = make(chan int)

	return cli
}

func (cli *CLI) SetToken(token string) {
	cli.token = token
}

func (cli *CLI) SetDay(day int) {
	cli.day = day
}

func (cli *CLI) ShowDetails(showDetails int) {
	cli.showDetails = showDetails
}

func (cli *CLI) Initialize() bool {
	user, _, err := cli.client.Users.Get(cli.context, "")
	if err != nil {
		fmt.Println("GitHub API authentication failed. Token may be invalid.", cli.token)
		return false
	}

	cli.username = string(user.GetLogin())
	return true
}

func (cli *CLI) GetRepos() []*github.Repository {
	// list all repositories for the authenticated user
	repos, _, err := cli.client.Repositories.List(cli.context, cli.username, nil)

	if err != nil {
		fmt.Println("Cannot fetch repositories from GitHub. Error message:", err)
		return []*github.Repository{}
	}

	if len(repos) == 0 {
		fmt.Println("User has no repositories")
		return []*github.Repository{}
	}

	return repos
}

func (cli *CLI) RepoStat(w io.Writer) {
	for repo := range cli.repoChannel {
		if cli.showDetails == 1 {
			fmt.Print(".")
		}

		stats, _, _ := cli.client.Repositories.ListTrafficViews(cli.context, cli.username, repo.GetName(), nil)

		viewCount := 0
		uniqueCount := 0
		dayDiff := int(math.Abs(float64(cli.day))) + 1

		for _, view := range stats.Views {
			if view.GetTimestamp().After(time.Now().UTC().AddDate(0, 0, -dayDiff)) {
				uniqueCount += view.GetUniques()
				viewCount += view.GetCount()
			}
		}

		if cli.showDetails == 1 {
			fmt.Fprintf(w, "%s\t%d\t%d\t\n", repo.GetFullName(), viewCount, uniqueCount)
		}

		cli.accountTotalChannel <- viewCount
		cli.accountUniqueChannel <- uniqueCount
	}

	cli.haltChannel <- 1
}

func (cli *CLI) Execute(w io.Writer) {
	fmt.Print("Checking repositories ")

	repos := cli.GetRepos()
	repoCount := len(repos)

	cli.accountTotalChannel = make(chan int, repoCount)
	cli.accountUniqueChannel = make(chan int, repoCount)

	if cli.showDetails == 1 {
		fmt.Fprintf(w, "Repository\tTotal View\tUnique View\t\n")
	}

	for i := 0; i < repoCount; i++ {
		go cli.RepoStat(w)
	}

	for _, repo := range repos {
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
			fmt.Println(".")
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

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Total View:", total, ", Unique View:", unique)
}
