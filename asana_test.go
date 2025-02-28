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
		name        string
		text        string
		workspaceID string
		projectID   string
		taskID      string
	}{
		{
			name:        `v0 full screen URL`,
			text:        `Here is the Asana task URL that you should read before doing code review.\r\nhttps://app.asana.com/0/364167036366785/1162650948650897/f\r\n`,
			workspaceID: ``,
			projectID:   `364167036366785`,
			taskID:      `1162650948650897`,
		},
		{
			name:        `v0 no full screen URL`,
			text:        `Task URL: https://app.asana.com/0/364167036366785/1162650948650897`,
			workspaceID: ``,
			projectID:   `364167036366785`,
			taskID:      `1162650948650897`,
		},
		{
			name:        `v1 URL without projectID`,
			text:        `Here is the Asana task URL that you should read before doing code review.\r\nhttps://app.asana.com/1/5590853349337/task/1209347772587937?focus=true\r\n`,
			workspaceID: `5590853349337`,
			projectID:   ``,
			taskID:      `1209347772587937`,
		},
		{
			name:        `v1 URL without projectID no query`,
			text:        `Here is the Asana task URL that you should read before doing code review.\r\nhttps://app.asana.com/1/5590853215187/task/1162650948650897\r\n`,
			workspaceID: `5590853215187`,
			projectID:   ``,
			taskID:      `1162650948650897`,
		},
		{
			name:        `v1 URL with projectID`,
			text:        `Task URL: https://app.asana.com/1/5590853215187/project/364167036366785/task/1162650948650897?focus=true`,
			workspaceID: `5590853215187`,
			projectID:   ``,
			taskID:      `1162650948650897`,
		},
		{
			name:        `No task URL`,
			text:        `Lorem ipsum dolor sit amet, consectetur adipiscing elit.`,
			workspaceID: ``,
			projectID:   ``,
			taskID:      ``,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			workspaceID, projectID, taskID := parseAsanaTaskLink(test.text)
			assert.Equal(t, test.workspaceID, workspaceID)
			assert.Equal(t, test.projectID, projectID)
			assert.Equal(t, test.taskID, taskID)
		})
	}
}

func TestAddPullRequestCommentToTask(t *testing.T) {
	c := asana.NewClientWithAccessToken(asanaToken)

	pr := loadRequestReviewRequestedEvent()

	taskID := createTask()

	_, err := AddPullRequestCommentToTask(c, taskID, requester, pr)
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

	pr := loadRequestReviewRequestedEvent()

	_, err := UpdateTaskComment(c, storyID, requester, pr)
	require.NoError(t, err)
}

func TestCodeReviewSubtask(t *testing.T) {
	c := asana.NewClientWithAccessToken(asanaToken)

	var subtask *asana.Task

	t.Run("AddCodeReviewSubtask", func(t *testing.T) {
		pr := loadRequestReviewRequestedEvent()

		taskID := createTask()
		due := asana.Date(time.Now().AddDate(0, 0, 3))

		var err error
		subtask, err = AddCodeReviewSubtask(c, taskID, 123, requester, reviewer, due, pr)
		require.NoError(t, err)
	})

	t.Run("AddCodeReviewSubtaskComment", func(t *testing.T) {
		pr := loadRequestReviewSubmittedEvent()

		_, err := AddCodeReviewSubtaskComment(c, subtask, requester, reviewer, pr)
		require.NoError(t, err)
	})
}

func TestFindSubtaskByName(t *testing.T) {
	c := asana.NewClientWithAccessToken(asanaToken)

	pr := loadRequestReviewRequestedEvent()

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
