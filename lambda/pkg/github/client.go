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
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewClient(ctx context.Context, token string) *Client {
	this := new(Client)
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
	if github := github.NewClient(tc); github == nil {
		return nil
	} else {
		this.Client = github
	}
	return this
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Client) ListRepos(ctx context.Context) (interface{}, error) {
	repos, _, err := this.Repositories.List(ctx, "", &github.RepositoryListOptions{})
	return repos, err
}
