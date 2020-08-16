package githubworkflow

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

const (
	baseGithubDomain = "github.com"
)

type Client interface {
	CreateBranch(context.Context, *CreateBranchRequest) (*CreateBranchResponse, error)
	CreatePullRequest(context.Context, *CreatePullRequestRequest) (*CreatePullRequestResponse, error)
}

type CreateBranchRequest struct {
	Owner         string
	Repo          string
	Reference     string
	Branch        string
	CommitMessage string
	AuthorName    string
	AuthorEmail   string

	// Changes are changes to add to the branch. It's map of
	// filename and its contents.
	Changes map[string]io.Reader
}

type CreateBranchResponse struct{}

type CreatePullRequestRequest struct {
	Owner  string
	Repo   string
	Base   string
	Branch string
	Title  string
	Body   string
	Draft  bool
}

type CreatePullRequestResponse struct {
	Number  int
	HTMLURL string
}

type client struct {
	githubClient        *github.Client
	gitClientAuthMethod transport.AuthMethod
}

func New(githubClient *github.Client, gitClientAuthMethod transport.AuthMethod) (Client, error) {
	return &client{
		githubClient:        githubClient,
		gitClientAuthMethod: gitClientAuthMethod,
	}, nil
}

func (c *client) CreateBranch(ctx context.Context, req *CreateBranchRequest) (*CreateBranchResponse, error) {
	gitRepository, err := git.CloneContext(ctx, memory.NewStorage(), memfs.New(), &git.CloneOptions{
		Depth:         1,
		URL:           fmt.Sprintf("https://%s/%s/%s", baseGithubDomain, req.Owner, req.Repo),
		ReferenceName: plumbing.NewBranchReferenceName(req.Reference),
		SingleBranch:  true,
		Auth:          c.gitClientAuthMethod,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %s", err)
	}

	worktree, err := gitRepository.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %s", err)
	}

	if err := worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(req.Branch),
		Create: true,
	}); err != nil {
		return nil, fmt.Errorf("failed to checkout to %s: %s", req.Branch, err)
	}

	for filename, contents := range req.Changes {
		f, err := worktree.Filesystem.Create(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to create file %s: %s", filename, err)
		}

		if _, err := io.Copy(f, contents); err != nil {
			return nil, fmt.Errorf("failed to write: %s", err)
		}
	}

	if err := worktree.AddGlob("."); err != nil {
		return nil, fmt.Errorf("failed to add changes to the index: %s", err)
	}

	if _, err := worktree.Commit(req.CommitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  req.AuthorName,
			Email: req.AuthorEmail,
			When:  time.Now(),
		},
	}); err != nil {
		return nil, fmt.Errorf("failed to commit changes: %s", err)
	}

	if err := gitRepository.PushContext(ctx, &git.PushOptions{
		Auth: c.gitClientAuthMethod,
	}); err != nil {
		return nil, fmt.Errorf("failed to push changes: %s", err)
	}

	return &CreateBranchResponse{}, nil
}

func (c *client) CreatePullRequest(ctx context.Context, req *CreatePullRequestRequest) (*CreatePullRequestResponse, error) {
	pr, _, err := c.githubClient.PullRequests.Create(ctx, req.Owner, req.Repo, &github.NewPullRequest{
		Title:               github.String(req.Title),
		Head:                github.String(req.Branch),
		Base:                github.String(req.Base),
		Body:                github.String(req.Body),
		MaintainerCanModify: github.Bool(true),
		Draft:               github.Bool(req.Draft),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create PR: %s", err)
	}
	return &CreatePullRequestResponse{
		Number:  pr.GetNumber(),
		HTMLURL: pr.GetHTMLURL(),
	}, nil
}
