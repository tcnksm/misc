package main

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	"github.com/tcnksm/misc/pkg/go-githubworkflow"
)

const (
	envGithubUser  = "GITHUB_USER"
	envGithubToken = "GITHUB_TOKEN"
)

func main() {
	gitClientAuthMethod := &http.BasicAuth{
		Username: os.Getenv(envGithubUser),
		Password: os.Getenv(envGithubToken),
	}

	githubClient := github.NewClient(oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: os.Getenv(envGithubToken),
		},
	)))

	client, err := githubworkflow.New(githubClient, gitClientAuthMethod)
	if err != nil {
		log.Fatal(err)
	}

	var (
		owner  = "tcnksm"
		repo   = "misc"
		branch = "test-branch"
	)

	ctx := context.Background()
	if _, err := client.CreateBranch(ctx, &githubworkflow.CreateBranchRequest{
		Owner:         owner,
		Repo:          repo,
		Reference:     "master",
		Branch:        branch,
		CommitMessage: "This is test commit message",
		AuthorName:    "go-githubworkflow",
		AuthorEmail:   "go-githubworkflow@github.com",
		Changes: map[string]io.Reader{
			"README.md": strings.NewReader(`This is modified by go-githubworkflow`),
		},
	}); err != nil {
		log.Fatal(err)
	}

	pr, err := client.CreatePullRequest(ctx, &githubworkflow.CreatePullRequestRequest{
		Owner:  owner,
		Repo:   repo,
		Base:   "master",
		Branch: branch,
		Title:  "This is test PR",
		Body:   `This PR is created by go-githubworkflow`,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("PR is created at", pr.HTMLURL)
}
