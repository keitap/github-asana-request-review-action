package githubasana

import (
	"testing"

	"bitbucket.org/mikehouston/asana-go"
	"github.com/google/go-github/v35/github"
	"github.com/stretchr/testify/require"
)

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

	h := NewHandler(conf, ac)

	pr := &github.PullRequestEvent{
		Action: pString("review_requested"),
		PullRequest: &github.PullRequest{
			Number:       pInt(1),
			Title:        pString("title"),
			Body:         pString("task is here:\nhttps://app.asana.com/0/1200243266984258/1200265547631636/f"),
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

	err := h.handlePullRequestEvent(pr)
	require.NoError(t, err)
}
