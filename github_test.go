package githubasana

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/google/go-github/v35/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadTestdata(filepath string) (name string, payload []byte) {
	s := strings.SplitN(path.Base(filepath), "-", 2)

	payload, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	return s[0], payload
}

func loadRequestReviewRequestedEvent() (*github.PullRequestEvent, error) {
	name, payload := loadTestdata("testdata/pull_request-review_requested.json")

	event, err := github.ParseWebHook(name, payload)
	if err != nil {
		panic(err)
	}

	return event.(*github.PullRequestEvent), nil
}

func loadRequestReviewSubmittedEvent() (*github.PullRequestReviewEvent, error) {
	name, payload := loadTestdata("testdata/pull_request_review-submitted-approved.json")

	event, err := github.ParseWebHook(name, payload)
	if err != nil {
		panic(err)
	}

	return event.(*github.PullRequestReviewEvent), nil
}

func TestParseRequestReviewerEvent(t *testing.T) {
	e, err := loadRequestReviewRequestedEvent()
	require.NoError(t, err)

	assert.Equal(t, "keitap-2nd", e.PullRequest.RequestedReviewers[0].GetLogin())
}
