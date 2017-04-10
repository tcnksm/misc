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
	"log"
	"os"
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

	org, repo, number := os.Args[1], os.Args[2], os.Args[3]
	body := strings.Join(os.Args[4:], " ")
	log.Printf("[INFO] Create a comment %q on https://github.com/%s/%s/pull/%s",
		body, org, repo, number,
	)

	// Construct github HTTP client
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	_ = client
}
