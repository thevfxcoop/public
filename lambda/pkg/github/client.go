package github

import (
	"context"

	// Modules
	github "github.com/google/go-github/github"
	oauth2 "golang.org/x/oauth2"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Client struct {
	*github.Client

	owner string
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewClient(ctx context.Context, token, owner string) *Client {
	this := new(Client)
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
	if github := github.NewClient(tc); github == nil {
		return nil
	} else {
		this.Client = github
		this.owner = owner
	}

	// Return success
	return this
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Client) ListRepos(ctx context.Context) (interface{}, error) {
	repos, _, err := this.Repositories.List(ctx, this.owner, &github.RepositoryListOptions{})
	return repos, err
}

func (this *Client) CreateIssue(ctx context.Context, repo, title, body string) (interface{}, error) {
	issue, _, err := this.Issues.Create(ctx, this.owner, repo, &github.IssueRequest{
		Title: &title,
		Body:  &body,
	})
	return issue, err
}
