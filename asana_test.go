package githubasana

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"bitbucket.org/mikehouston/asana-go"
	"github.com/google/go-github/v74/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	projectID                 = "1200261405938356"
	storyID                   = "1209537017377484"
	emptyTaskID               = "1209536608330908"
	hasSignatureCommentTaskID = "1209536608330906"
	hasSubtaskTaskID          = "1209536608330912"
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

// requireAsanaToken guards tests that hit the live Asana API. It skips only
// when SKIP_INTEGRATION_TEST is explicitly set (e.g. Dependabot runs, which
// have no access to secrets); otherwise a missing token fails loudly so a
// misconfigured secret is not silently treated as a pass.
func requireAsanaToken(t *testing.T) {
	t.Helper()

	if os.Getenv("SKIP_INTEGRATION_TEST") == "true" {
		t.Skip("SKIP_INTEGRATION_TEST is set; skipping Asana integration test")
	}

	if asanaToken == "" {
		t.Fatal("ASANA_TOKEN is not set (set SKIP_INTEGRATION_TEST=true to skip integration tests)")
	}
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

func commentHTML(status, suffix string) string {
	reviewURL := "https://github.com/owner/repo/pull/1#pullrequestreview-123"
	reviewerPermalink := `<a data-asana-gid="12345"/>`

	return `<body><a href="` + reviewURL + `"><b>` + status + `</b></a>: reviewed by ` + reviewerPermalink + suffix + `</body>`
}

func TestBuildReviewCommentHTML(t *testing.T) {
	reviewURL := "https://github.com/owner/repo/pull/1#pullrequestreview-123"
	reviewerPermalink := `<a data-asana-gid="12345"/>`

	tests := []struct {
		name         string
		state, body  string
		commentCount int
		expected     string
	}{
		{"approved without comments", "approved", "", 0, commentHTML("✅ Approved", "")},
		{"approved with multiple comments", "approved", "", 3, commentHTML("✅💬 Approved with 3 comments", "")},
		{"approved with single comment", "approved", "", 1, commentHTML("✅💬 Approved with 1 comment", "")},
		{"approved with body and comments", "approved", "LGTM!", 3, commentHTML("✅💬 Approved with 3 comments", "\nLGTM!")},
		{"approved with body only", "approved", "LGTM!", 0, commentHTML("✅ Approved", "\nLGTM!")},
		{"changes_requested with comments", "changes_requested", "Please fix", 2, commentHTML("❗️💬 Changes Requested with 2 comments", "\nPlease fix")},
		{"commented with single comment", "commented", "", 1, commentHTML("💬 Commented with 1 comment", "")},
		{"commented with multiple comments", "commented", "", 5, commentHTML("💬 Commented with 5 comments", "")},
		{"body with whitespace is trimmed", "approved", "  LGTM!  ", 0, commentHTML("✅ Approved", "\nLGTM!")},
		{"body with markdown blockquote is escaped", "approved", "> quoted", 0, commentHTML("✅ Approved", "\n&gt; quoted")},
		{"body with html special chars is escaped", "approved", `a <b> & "c"`, 0, commentHTML("✅ Approved", "\na &lt;b&gt; &amp; &#34;c&#34;")},
		{
			"approved with comments and multiline blockquote body",
			"approved", "> quoted line from a previous comment\n\nPlease keep a record of the source in the PR.", 3,
			commentHTML("✅💬 Approved with 3 comments", "\n&gt; quoted line from a previous comment\n\nPlease keep a record of the source in the PR."),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := buildReviewCommentHTML(reviewURL, test.state, reviewerPermalink, test.body, test.commentCount)
			assert.Equal(t, test.expected, got)
		})
	}
}

func TestAddPullRequestCommentToTask(t *testing.T) {
	requireAsanaToken(t)

	c := asana.NewClientWithAccessToken(asanaToken)

	pr := loadRequestReviewRequestedEvent()

	taskID := createTask()

	_, err := AddPullRequestCommentToTask(c, taskID, requester, pr)
	require.NoError(t, err)
}

func TestFindTaskComment(t *testing.T) {
	requireAsanaToken(t)

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

// newFakeAsanaClient returns a client pointed at a local fake Asana API.
func newFakeAsanaClient(t *testing.T, handler http.Handler) *asana.Client {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	c := asana.NewClient(srv.Client())

	u, err := url.Parse(srv.URL)
	require.NoError(t, err)

	c.BaseURL = u

	return c
}

func writeAsanaPage(t *testing.T, w http.ResponseWriter, data any, nextOffset string) {
	t.Helper()

	resp := map[string]any{"data": data}
	if nextOffset != "" {
		resp["next_page"] = map[string]string{"offset": nextOffset}
	}

	w.Header().Set("Content-Type", "application/json")
	require.NoError(t, json.NewEncoder(w).Encode(resp))
}

// TestFindTaskCommentAcrossPages reproduces the duplicate comment incident:
// the stories endpoint returns entries oldest first, so on a task with more
// than 100 stories the bot's comment only appears on a later page.
func TestFindTaskCommentAcrossPages(t *testing.T) {
	fillerPage := make([]map[string]string, 100)
	for i := range fillerPage {
		fillerPage[i] = map[string]string{"gid": fmt.Sprintf("sys-%d", i), "text": "changed the due date"}
	}

	tests := []struct {
		name       string
		secondPage []map[string]string
		expectedID string
	}{
		{
			name: "comment on second page is found",
			secondPage: []map[string]string{
				{"gid": "story-101", "text": "PR comment\n\n" + signature},
			},
			expectedID: "story-101",
		},
		{
			name:       "no comment on any page",
			secondPage: []map[string]string{{"gid": "sys-101", "text": "liked the task"}},
			expectedID: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/tasks/42/stories", func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Query().Get("offset") {
				case "":
					writeAsanaPage(t, w, fillerPage, "page2")
				case "page2":
					writeAsanaPage(t, w, test.secondPage, "")
				default:
					http.Error(w, "unexpected offset", http.StatusBadRequest)
				}
			})

			c := newFakeAsanaClient(t, mux)

			story, err := FindTaskComment(c, "42", signature)
			require.NoError(t, err)

			if test.expectedID == "" {
				assert.Nil(t, story)
			} else {
				require.NotNil(t, story)
				assert.Equal(t, test.expectedID, story.ID)
			}
		})
	}
}

func TestFindSubtaskByNameAcrossPages(t *testing.T) {
	fillerPage := make([]map[string]string, 100)
	for i := range fillerPage {
		fillerPage[i] = map[string]string{"gid": fmt.Sprintf("sub-%d", i), "name": fmt.Sprintf("subtask %d", i)}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks/42/subtasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("offset") {
		case "":
			writeAsanaPage(t, w, fillerPage, "page2")
		case "page2":
			writeAsanaPage(t, w, []map[string]string{
				{"gid": "sub-101", "name": "✍️ Code review: #123 Keita Kitamura (@keitap)"},
			}, "")
		default:
			http.Error(w, "unexpected offset", http.StatusBadRequest)
		}
	})

	c := newFakeAsanaClient(t, mux)

	subtask, err := FindSubtaskByName(c, "42", "@keitap")
	require.NoError(t, err)
	require.NotNil(t, subtask)
	assert.Equal(t, "sub-101", subtask.ID)
}

func TestUpdateTaskComment(t *testing.T) {
	requireAsanaToken(t)

	c := asana.NewClientWithAccessToken(asanaToken)

	pr := loadRequestReviewRequestedEvent()

	_, err := UpdateTaskComment(c, storyID, requester, pr)
	require.NoError(t, err)
}

func TestCodeReviewSubtask(t *testing.T) {
	requireAsanaToken(t)

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

		_, err := AddCodeReviewSubtaskComment(c, subtask, requester, reviewer, pr, 3)
		require.NoError(t, err)
	})

	t.Run("AddCodeReviewSubtaskComment with blockquote body", func(t *testing.T) {
		pr := loadRequestReviewSubmittedEvent()

		// A review body starting with a markdown blockquote previously made the
		// Asana html_text invalid and rendered the whole comment as raw HTML.
		// Posting it to Asana verifies the escaped html_text renders correctly.
		pr.Review.Body = github.Ptr("> quoted line from a previous comment\n\nPlease keep a record of the source in the PR.")

		_, err := AddCodeReviewSubtaskComment(c, subtask, requester, reviewer, pr, 3)
		require.NoError(t, err)
	})
}

func TestFindSubtaskByName(t *testing.T) {
	requireAsanaToken(t)

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
