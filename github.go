package githubasana

import (
	"context"
	"fmt"

	"github.com/google/go-github/v71/github"
)

func getRequestedReviewers(gh *github.Client, owner string, repo string, number int) ([]*github.User, error) {
	reviewers, _, err := gh.PullRequests.ListReviewers(context.Background(), owner, repo, number, &github.ListOptions{PerPage: 100})
	if err != nil {
		return nil, fmt.Errorf("failed to get reviewers: %w", err)
	}

	return reviewers.Users, nil
}
