package githubasana

import (
	"errors"

	"github.com/google/go-github/v35/github"
)

func parseRequestReviewerEvent(name string, payload []byte) (*github.PullRequestEvent, error) {
	event, err := github.ParseWebHook(name, payload)
	if err != nil {
		return nil, err
	}

	switch event := event.(type) {
	case *github.PullRequestEvent:
		return event, nil
	}

	return nil, errors.New("unknown event type: " + name)
}
