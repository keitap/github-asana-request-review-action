package githubasana

import (
	"testing"

	"bitbucket.org/mikehouston/asana-go"
	"github.com/google/go-github/v35/github"
	"github.com/stretchr/testify/require"
)

var prEventReviewRequested = &github.PullRequestEvent{
	Action: pString("review_requested"),
	PullRequest: &github.PullRequest{
		State:        pString("open"),
		Number:       pInt(1),
		Title:        pString("title"),
		Body:         pString("task is here:\nhttps://app.asana.com/0/1200261405938356/1209536608330915"),
		ChangedFiles: pInt(1),
		Additions:    pInt(2),
		Deletions:    pInt(3),
		HTMLURL:      pString("https://github.com/keitap/github-actions-test/pull/1"),
		RequestedReviewers: []*github.User{
			{
				Login: pString("keitap-2nd"),
			},
			{
				Login: pString("no-asana-user"),
			},
		},
		User: &github.User{
			Login: pString("keitap"),
		},
	},
}

func pString(s string) *string {
	return &s
}

func pInt(i int) *int {
	return &i
}

func TestHandler_handlePullRequestEvent(t *testing.T) {
	conf := &Config{
		Accounts: map[GithubLogin]AsanaGID{
			"keitap":     "5590853215184",
			"keitap-2nd": "2540808972045",
		},
	}

	ac := asana.NewClientWithAccessToken(asanaToken)
	gh := github.NewClient(nil)

	h := NewHandler(conf, ac, gh)

	err := h.handlePullRequestEvent(prEventReviewRequested)
	require.NoError(t, err)
}

func TestHandler_handlePullRequestEvent_NoConfig(t *testing.T) {
	conf := &Config{}

	ac := asana.NewClientWithAccessToken(asanaToken)
	gh := github.NewClient(nil)

	h := NewHandler(conf, ac, gh)

	err := h.handlePullRequestEvent(prEventReviewRequested)
	require.NoError(t, err)
}
