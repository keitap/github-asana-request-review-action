package pkg

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"bitbucket.org/mikehouston/asana-go"
	"github.com/google/go-github/v35/github"
)

const signature = "#github-asana-request-review-action"

var (
	taskURLMatcher = regexp.MustCompile(`https://app.asana.com/0/(\d+)/(\d+)`)
)

func parseAsanaTaskLink(text string) (projectID string, taskID string) {
	m := taskURLMatcher.FindStringSubmatch(text)
	if len(m) <= 0 {
		return "", ""
	}

	return m[1], m[2]
}

func AddPullRequestURLToTaskComment(client *asana.Client, taskID string, pr *github.PullRequestEvent) (*asana.Story, error) {
	task := &asana.Task{ID: taskID}
	story := &asana.StoryBase{
		HTMLText: createPRText(pr),
		IsPinned: true,
	}

	return task.CreateComment(client, story)
}

// FindTaskComment finds a story which contains specified string.
func FindTaskComment(client *asana.Client, taskID string, findString string) (*asana.Story, error) {
	task := &asana.Task{ID: taskID}
	stories, _, err := task.Stories(client, &asana.Options{
		Limit: 100,
	})
	if err != nil {
		return nil, err
	}

	for _, s := range stories {
		if strings.Contains(s.Text, findString) {
			return s, nil
		}
	}

	return nil, nil
}

func UpdateTaskComment(client *asana.Client, storyID string, pr *github.PullRequestEvent) (*asana.Story, error) {
	story := &asana.Story{ID: storyID}
	newStory := &asana.StoryBase{
		HTMLText: createPRText(pr),
		IsPinned: true,
	}

	return story.UpdateStory(client, newStory)
}

func AddCodeReviewSubtask(client *asana.Client, taskID string, assigneeUserID string, dueDate asana.Date, pr *github.PullRequestEvent) (*asana.Task, error) {
	req := &asana.CreateTaskRequest{
		Parent:   taskID,
		Assignee: assigneeUserID,
		TaskBase: asana.TaskBase{
			Name:      fmt.Sprintf(`Code Review request to %s`, pr.RequestedReviewer.GetLogin()),
			HTMLNotes: createPRText(pr),
			DueOn:     &dueDate,
		},
	}

	return client.CreateTask(req)
}

// FindSubtaskByName finds a subtask which contains specified string.
func FindSubtaskByName(client *asana.Client, taskID string, findString string) (*asana.Task, error) {
	task := &asana.Task{ID: taskID}

	subtasks, _, err := task.Subtasks(client, &asana.Options{
		Limit: 100,
	})
	if err != nil {
		return nil, err
	}

	for _, s := range subtasks {
		if strings.Contains(s.Name, findString) {
			return s, nil
		}
	}

	return nil, nil
}

func createPRText(pr *github.PullRequestEvent) string {
	return fmt.Sprintf(`<body>Pull request is created by <b>%s</b>.
<a href="%s">#%d: %s</a>

<b>%d</b> changed files (<b>+%d -%d</b>)

%s
</body>`,
		pr.Sender.GetLogin(),
		pr.PullRequest.GetHTMLURL(), pr.PullRequest.GetNumber(), pr.PullRequest.GetTitle(),
		pr.PullRequest.GetChangedFiles(), pr.PullRequest.GetAdditions(), pr.PullRequest.GetDeletions(),
		signature+` at `+time.Now().Format(time.RFC3339),
	)
}
