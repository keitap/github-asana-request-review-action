package githubasana

import (
	"fmt"
	"regexp"
	"strings"

	"bitbucket.org/mikehouston/asana-go"
	"github.com/google/go-github/v35/github"
)

const signature = "#github-asana-request-review"

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

func AddPullRequestCommentToTask(client *asana.Client, taskID string, pr *github.PullRequestEvent) (*asana.Story, error) {
	task := &asana.Task{ID: taskID}
	story := &asana.StoryBase{
		HTMLText: createPullRequestCommentText(pr),
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
		HTMLText: createPullRequestCommentText(pr),
		IsPinned: true,
	}

	return story.UpdateStory(client, newStory)
}

func AddCodeReviewSubtask(client *asana.Client, taskID string, requesterUserID string, assigneeUserID string, reviewer string, dueDate asana.Date, pr *github.PullRequestEvent) (*asana.Task, error) {
	req := &asana.CreateTaskRequest{
		Parent:    taskID,
		Assignee:  assigneeUserID,
		Followers: []string{requesterUserID},
		TaskBase: asana.TaskBase{
			Name:      fmt.Sprintf(`Code review request to %s`, reviewer),
			HTMLNotes: createReviewRequestDescText(pr),
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

func createPullRequestCommentText(pr *github.PullRequestEvent) string {
	reviewers := make([]string, len(pr.PullRequest.RequestedReviewers))
	for i, u := range pr.PullRequest.RequestedReviewers {
		reviewers[i] = fmt.Sprintf("<b>%s</b>", u.GetLogin())
	}

	return fmt.Sprintf(`<body><code><a href="%s">Pull request #%d: %s</a>

<b>%d</b> changed files (<b>+%d -%d</b>) created by <b>%s</b>
Reviewers: %s

by %s
</code></body>`,
		pr.PullRequest.GetHTMLURL(), pr.PullRequest.GetNumber(), pr.PullRequest.GetTitle(),
		pr.PullRequest.GetChangedFiles(), pr.PullRequest.GetAdditions(), pr.PullRequest.GetDeletions(), pr.Sender.GetLogin(),
		strings.Join(reviewers, ", "),
		signature,
	)
}

func createReviewRequestDescText(pr *github.PullRequestEvent) string {
	return fmt.Sprintf(`<body><code>Could you please review a pull request created by <b>%s</b> 🙇
<a href="%s">#%d: %s</a>

After you finished a code review, pass this assign back to <b>%s</b>.
Do not mark complete unless you are <b>%s</b>.
</code></body>`,
		pr.Sender.GetLogin(),
		pr.PullRequest.GetHTMLURL(), pr.PullRequest.GetNumber(), pr.PullRequest.GetTitle(),
		pr.Sender.GetLogin(),
		pr.Sender.GetLogin(),
	)
}
