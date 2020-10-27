/*
Command 'github-label' add label on GitHub issues and PR

  $ github-label ORG REPO NUMBER LABEL...

For example, if you want to add label "team/x" on https://github.com/tcnksm/ghr pull request number 987

  $ github-label tcnksm ghr 987 "team/x"

To use this command, you need to prepare GitHub API Token and set it via GITHUB_TOKEN
env var.

To install, use go get,

  $ go get github.com/tcnksm/misc/cmd/github-label

*/
package main

import (
	"log"
	"os"
	"strconv"

	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const EnvToken = "GITHUB_TOKEN"

func main() {

	if len(os.Args) < 5 {
		log.Fatal("[Usage] github-label ORG REPO NUMBER LABEL...")
	}

	token := os.Getenv(EnvToken)
	if len(token) == 0 {
		log.Fatal("You need GitHub API token via GITHUB_TOKEN env var")
	}

	var (
		issueNumber int
		err         error
	)
	owner, repo, issueNumberStr := os.Args[1], os.Args[2], os.Args[3]
	labels := os.Args[4:]

	issueNumber, err = strconv.Atoi(issueNumberStr)
	if err != nil {
		log.Fatalf("[ERROR] Issue number must be int: %e", err)
	}

	// Construct github HTTP client
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	log.Printf("[INFO] Add labels on https://github.com/%s/%s/issues/%d", owner, repo, issueNumber)
	if _, _, err := client.Issues.AddLabelsToIssue(context.Background(), owner, repo, issueNumber, labels); err != nil {
		log.Fatalf("[ERROR] Failed to add labels: %s", err)
	}

	log.Printf("[INFO] Successfully added labels!")
}
