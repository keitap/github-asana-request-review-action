package pkg

import (
	"os"
	"testing"
	"time"

	"bitbucket.org/mikehouston/asana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	AssigneeUserID        = "5590853215184"
	TaskID                = "1200243266984261"
	HasPRURLCommentTaskID = "1200243266984263"
	NoPRURLCommentTaskID  = "1200243266984265"
)

var (
	asanaAccessToken = ""
)

func init() {
	asanaAccessToken = os.Getenv("ASANA_ACCESS_TOKEN")
}

func TestParseAsanaTaskLink(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		projectID string
		taskID    string
	}{
		{
			name:      `full screen URL`,
			text:      `Here is the Asana task URL that you should read before doing code review.\r\nhttps://app.asana.com/0/364167036366785/1162650948650897/f\r\n`,
			projectID: `364167036366785`,
			taskID:    `1162650948650897`,
		},
		{
			name:      `no full screen URL`,
			text:      `Task URL: https://app.asana.com/0/364167036366785/1162650948650897`,
			projectID: `364167036366785`,
			taskID:    `1162650948650897`,
		},
		{
			name:      `No task URL`,
			text:      `Lorem ipsum dolor sit amet, consectetur adipiscing elit.`,
			projectID: ``,
			taskID:    ``,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			projectID, taskID := parseAsanaTaskLink(test.text)
			assert.Equal(t, test.projectID, projectID)
			assert.Equal(t, test.taskID, taskID)
		})
	}
}

func TestAddPullRequestURLToTaskComment(t *testing.T) {
	c := asana.NewClientWithAccessToken(asanaAccessToken)

	pr, err := loadRequestReviewerEvent()
	require.NoError(t, err)

	_, err = AddPullRequestURLToTaskComment(c, TaskID, pr)
	require.NoError(t, err)
}

func TestHasCommentContainsURL(t *testing.T) {
	c := asana.NewClientWithAccessToken(asanaAccessToken)

	pr, err := loadRequestReviewerEvent()
	require.NoError(t, err)

	url := *pr.PullRequest.HTMLURL

	tests := []struct {
		name     string
		taskID   string
		expected bool
	}{
		{name: "has comment", taskID: HasPRURLCommentTaskID, expected: true},
		{name: "no comment", taskID: NoPRURLCommentTaskID, expected: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, err := HasCommentContainsURL(c, test.taskID, url)
			require.NoError(t, err)
			assert.Equal(t, test.expected, value)
		})
	}
}

func TestAddCodeReviewSubtask(t *testing.T) {
	c := asana.NewClientWithAccessToken(asanaAccessToken)

	pr, err := loadRequestReviewerEvent()
	require.NoError(t, err)

	due := asana.Date(time.Now().AddDate(0, 0, 3))

	_, err = AddCodeReviewSubtask(c, TaskID, AssigneeUserID, due, pr)
	require.NoError(t, err)
}
