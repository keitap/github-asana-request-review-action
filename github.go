package githubasana

import (
	"context"
	"fmt"

	"github.com/google/go-github/v74/github"
)

func getRequestedReviewers(gh *github.Client, owner string, repo string, number int) ([]*github.User, error) {
	reviewers, _, err := gh.PullRequests.ListReviewers(context.Background(), owner, repo, number, &github.ListOptions{PerPage: 100})
	if err != nil {
		return nil, fmt.Errorf("failed to get reviewers: %w", err)
	}

	return reviewers.Users, nil
}

func getReviewCommentCount(gh *github.Client, owner string, repo string, number int, reviewID int64) (int, error) {
	comments, _, err := gh.PullRequests.ListReviewComments(context.Background(), owner, repo, number, reviewID, &github.ListOptions{PerPage: 100})
	if err != nil {
		return 0, fmt.Errorf("failed to get review comments: %w", err)
	}

	return len(comments), nil
}
