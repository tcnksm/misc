/*
Command 'github-pull-request' creates a PR on GitHub.

  $ github-pull-request [OPTIONS...] ORG REPO BRANCH

The difference from the Github official command https://cli.github.com/ is you can
specify additional information like the label or milestone or reviwers.

To use this command, you need to prepare GitHub API Token and set it via GITHUB_TOKEN
env var.

To install, use go get,

  $ go get github.com/tcnksm/misc/cmd/github-pull-request

*/
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
)

const EnvToken = "GITHUB_TOKEN"

var usage = `Usage: github-pull-requst [options...] ORG REPO BRANCH

Options:
  -base string    The branch into which you want your code merged. By default, master is used.
  -title string   The pull request title. By default, branch name is used.
  -body string    The pull request body. By default, it's empty.
  -labels string  The comma separated list of labels to put the PR
`

var (
	baseF  = flag.String("base", "master", "")
	titleF = flag.String("title", "", "")
	bodyF  = flag.String("body", "", "")
)

func main() {
	ctx := context.Background()

	var labels labels
	flag.Var(&labels, "labels", "")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 3 {
		log.Println("[ERROR] Invalid argument")
		fmt.Fprint(os.Stderr, usage)
		return
	}
	owner, repo, branch := args[0], args[1], args[2]

	base := *baseF
	title := *titleF
	body := *bodyF
	if title == "" {
		title = branch
	}

	token := os.Getenv(EnvToken)
	if len(token) == 0 {
		log.Fatal("You need GitHub API token via GITHUB_TOKEN env var")
	}

	// Construct github HTTP client
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	pr, _, err := client.PullRequests.Create(ctx, owner, repo, &github.NewPullRequest{
		Title: &title,
		Head:  &branch,
		Base:  &base,
		Body:  &body,
	})
	if err != nil {
		log.Fatalf("[ERROR] Faield to create a PR: %s", err)
	}
	log.Printf("[INFO] Successfully created a PR: %s", *pr.HTMLURL)

	if len(labels) > 0 {
		_, _, err := client.Issues.AddLabelsToIssue(ctx, owner, repo, *pr.Number, labels)
		if err != nil {
			log.Fatalf("[ERROR] Faield to add labels: %s", err)
		}
		log.Println("[INFO] Successfully labels are added")
	}
}

type labels []string

func (l *labels) String() string {
	return strings.Join(*l, ",")
}

func (l *labels) Set(s string) error {
	items := strings.Split(s, ",")
	for _, item := range items {
		if !contains(*l, item) {
			*l = append(*l, item)
		}
	}
	return nil
}

func contains(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
