/*
Command 'github-comment' comments on GitHub pull request.

  $ github-comment ORG REPO NUMBER COMMENT

For example, if you want to comment on https://github.com/tcnksm/ghr pull request number 987

  $ github-comment tcnksm ghr 987 "This is test comment"

To use this command, you need to prepare GitHub API Token and set it via GITHUB_TOKEN
env var.

To install, use go get,

  $ go get github.com/tcnksm/misc/cmd/github-comment

*/
package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const EnvToken = "GITHUB_TOKEN"

func main() {

	if len(os.Args) < 5 {
		log.Fatal("[Usage] github-comment ORG REPO NUMBER BODY")
	}

	token := os.Getenv(EnvToken)
	if len(token) == 0 {
		log.Fatal("You need GitHub API token via GITHUB_TOKEN env var")
	}

	owner, repo := os.Args[1], os.Args[2]
	issueNumber, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatalf("[ERROR] Issue number must be int: %e", err)
	}

	body := strings.Join(os.Args[4:], " ")
	log.Printf("[INFO] Create a comment %q on https://github.com/%s/%s/issues/%d",
		body, owner, repo, issueNumber,
	)

	// Construct github HTTP client
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	// Check there are no same comments
	comments, _, err := client.Issues.ListComments(context.Background(), owner, repo, issueNumber, nil)
	if err != nil {
		log.Fatalf("[ERROR] Failed to get comments: %s", err)
	}
	for _, c := range comments {
		if c.GetBody() == body {
			log.Printf("[INFO] comment %q was already posted, skip it", body)
			os.Exit(0)
		}
	}

	if _, _, err := client.Issues.CreateComment(context.Background(), owner, repo, issueNumber, &github.IssueComment{
		Body: &body,
	}); err != nil {
		log.Fatalf("[ERROR] Failed to create comment: %s", err)
	}

	log.Printf("[INFO] Successfully created a comment!")
}
