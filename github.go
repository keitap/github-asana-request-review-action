package githubasana

import (
	"context"
	"errors"

	"github.com/google/go-github/v35/github"
	"golang.org/x/xerrors"
)

func parseRequestReviewerEvent(name string, payload []byte) (*github.PullRequestEvent, error) {
	event, err := github.ParseWebHook(name, payload)
	if err != nil {
		return nil, err
	}

	switch event := event.(type) {
	case *github.PullRequestEvent:
		return event, nil
	default:
		return nil, errors.New("unknown event type: " + name)
	}
}

func getRequestedReviewers(gh *github.Client, owner string, repo string, number int) ([]*github.User, error) {
	reviewers, _, err := gh.PullRequests.ListReviewers(context.Background(), owner, repo, number, &github.ListOptions{PerPage: 100})
	if err != nil {
		return nil, xerrors.Errorf("failed to get reviewers: %w", err)
	}

	return reviewers.Users, nil
}
