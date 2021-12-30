package githubasana

import (
	"os"
	"testing"
	"time"

	"bitbucket.org/mikehouston/asana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	projectID                 = "1200243266984258"
	storyID                   = "1200243344965037"
	emptyTaskID               = "1200243266984265"
	hasSignatureCommentTaskID = "1200243266984263"
	hasSubtaskTaskID          = "1200243529563651"
)

var (
	asanaToken = ""
	taskID     = ""

	requester = &Account{
		Name:         "Keita Kitamura",
		AsanaUserGID: "5590853215184",
		GitHubLogin:  "keitap",
	}

	reviewer = &Account{
		Name:         "Keita Kitamura",
		AsanaUserGID: "2540808972045",
		GitHubLogin:  "keitap",
	}
)

func init() {
	asanaToken = os.Getenv("ASANA_TOKEN")
}

func createTask() string {
	if taskID != "" {
		return taskID
	}

	c := asana.NewClientWithAccessToken(asanaToken)

	req := &asana.CreateTaskRequest{
		TaskBase: asana.TaskBase{
			Name:  "Test " + time.Now().Format(time.RFC3339),
			Notes: "This task is for testing purpose.\nNo problem to delete as you like.",
		},
		Projects: []string{projectID},
	}

	task, err := c.CreateTask(req)
	if err != nil {
		panic(err)
	}

	taskID = task.ID

	return taskID
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

func TestAddPullRequestCommentToTask(t *testing.T) {
	c := asana.NewClientWithAccessToken(asanaToken)

	pr, err := loadRequestReviewerEvent()
	require.NoError(t, err)

	taskID := createTask()

	_, err = AddPullRequestCommentToTask(c, taskID, requester, pr)
	require.NoError(t, err)
}

func TestFindTaskComment(t *testing.T) {
	c := asana.NewClientWithAccessToken(asanaToken)

	tests := []struct {
		name     string
		taskID   string
		expected bool
	}{
		{name: "has comment", taskID: hasSignatureCommentTaskID, expected: true},
		{name: "no comment", taskID: emptyTaskID, expected: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			story, err := FindTaskComment(c, test.taskID, signature)
			require.NoError(t, err)
			assert.Equal(t, test.expected, story != nil)
		})
	}
}

func TestUpdateTaskComment(t *testing.T) {
	c := asana.NewClientWithAccessToken(asanaToken)

	pr, err := loadRequestReviewerEvent()
	require.NoError(t, err)

	_, err = UpdateTaskComment(c, storyID, requester, pr)
	require.NoError(t, err)
}

func TestAddCodeReviewSubtask(t *testing.T) {
	c := asana.NewClientWithAccessToken(asanaToken)

	pr, err := loadRequestReviewerEvent()
	require.NoError(t, err)

	taskID := createTask()
	due := asana.Date(time.Now().AddDate(0, 0, 3))

	_, err = AddCodeReviewSubtask(c, taskID, requester, reviewer, due, pr)
	require.NoError(t, err)
}

func TestFindSubtaskByName(t *testing.T) {
	c := asana.NewClientWithAccessToken(asanaToken)

	pr, err := loadRequestReviewerEvent()
	require.NoError(t, err)

	githubReviewerLogin := pr.PullRequest.RequestedReviewers[0].GetLogin()

	tests := []struct {
		name     string
		taskID   string
		expected bool
	}{
		{name: "has subtask", taskID: hasSubtaskTaskID, expected: true},
		{name: "no subtask", taskID: emptyTaskID, expected: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			subtask, err := FindSubtaskByName(c, test.taskID, githubReviewerLogin)
			require.NoError(t, err)
			assert.Equal(t, test.expected, subtask != nil)
		})
	}
}
