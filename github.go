package githubasana

import (
	"context"

	"github.com/google/go-github/v35/github"
	"golang.org/x/xerrors"
)

func getRequestedReviewers(gh *github.Client, owner string, repo string, number int) ([]*github.User, error) {
	reviewers, _, err := gh.PullRequests.ListReviewers(context.Background(), owner, repo, number, &github.ListOptions{PerPage: 100})
	if err != nil {
		return nil, xerrors.Errorf("failed to get reviewers: %w", err)
	}

	return reviewers.Users, nil
}
