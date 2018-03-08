/*
Command 'github-comment' comments on GitHub issues/revision.

  $ github-comment ORG REPO NUMBER|REVISION COMMENT

For example, if you want to comment on https://github.com/tcnksm/ghr pull request number 987

  $ github-comment tcnksm ghr 987 "This is test comment"

If you want to do on https://github.com/tcnksm/ghr/commit/5e0f14a4236b60e9c59aa3523cd7dac60e1859e8

  $ github-comment tcnksm ghr 5e0f14a "This is test comment too"

If you comment on a revision, hash must be 7 or more chars

To use this command, you need to prepare GitHub API Token and set it via GITHUB_TOKEN
env var.

To install, use go get,

  $ go get github.com/tcnksm/misc/cmd/github-comment

*/
package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const EnvToken = "GITHUB_TOKEN"

func main() {

	if len(os.Args) < 5 {
		log.Fatal("[Usage] github-comment ORG REPO NUMBER|REVISION BODY")
	}

	token := os.Getenv(EnvToken)
	if len(token) == 0 {
		log.Fatal("You need GitHub API token via GITHUB_TOKEN env var")
	}

	var (
		issueNumber int
		err         error
	)
	owner, repo, revision := os.Args[1], os.Args[2], os.Args[3]

	isRevision := isValidRevision(revision)
	if !isRevision {
		issueNumber, err = strconv.Atoi(revision)
		if err != nil {
			log.Fatalf("[ERROR] Issue number must be int: %e", err)
		}
	}

	body := strings.Join(os.Args[4:], " ")

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

	if isRevision {
		log.Printf("[INFO] Creating a comment %q on https://github.com/%s/%s/commit/%s",
			body, owner, repo, revision,
		)
		_, _, err = client.Repositories.CreateComment(
			context.Background(),
			owner,
			repo,
			revision,
			&github.RepositoryComment{Body: &body},
		)
	} else {
		log.Printf("[INFO] Creating a comment %q on https://github.com/%s/%s/issues/%d",
			body, owner, repo, issueNumber,
		)
		_, _, err = client.Issues.CreateComment(
			context.Background(),
			owner,
			repo,
			issueNumber,
			&github.IssueComment{Body: &body},
		)
	}
	if err != nil {
		log.Fatalf("[ERROR] Failed to create comment: %s", err)
	}

	log.Printf("[INFO] Successfully created a comment!")
}

func isValidRevision(arg string) bool {
	// In GitHub, 7 chars word is used as a short hash
	if len(arg) >= 7 {
		return true
	}
	return false
}
