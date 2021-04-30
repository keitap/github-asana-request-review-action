package githubasana

import (
	"fmt"
	"regexp"
	"strings"

	"bitbucket.org/mikehouston/asana-go"
	"github.com/google/go-github/v35/github"
)

const signature = "#github-asana-request-review"

var taskURLMatcher = regexp.MustCompile(`https://app.asana.com/0/(\d+)/(\d+)`)

func parseAsanaTaskLink(text string) (projectID string, taskID string) {
	m := taskURLMatcher.FindStringSubmatch(text)
	if len(m) == 0 {
		return "", ""
	}

	return m[1], m[2]
}

func AddPullRequestCommentToTask(client *asana.Client, taskID string, requester *Account, reviewers []*Account, pr *github.PullRequestEvent) (*asana.Story, error) {
	task := &asana.Task{ID: taskID}
	story := &asana.StoryBase{
		HTMLText: createPullRequestCommentText(requester, reviewers, pr),
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

func UpdateTaskComment(client *asana.Client, storyID string, requester *Account, reviewers []*Account, pr *github.PullRequestEvent) (*asana.Story, error) {
	story := &asana.Story{ID: storyID}
	newStory := &asana.StoryBase{
		HTMLText: createPullRequestCommentText(requester, reviewers, pr),
		IsPinned: true,
	}

	return story.UpdateStory(client, newStory)
}

func AddCodeReviewSubtask(client *asana.Client, taskID string, requester *Account, reviewer *Account, dueDate asana.Date, pr *github.PullRequestEvent) (*asana.Task, error) {
	req := &asana.CreateTaskRequest{
		Parent:    taskID,
		Assignee:  reviewer.AsanaUserGID,
		Followers: []string{requester.AsanaUserGID},
		TaskBase: asana.TaskBase{
			Name:      fmt.Sprintf(`✍️ Code review: %s`, reviewer),
			HTMLNotes: createReviewRequestDescText(requester, pr),
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

func createPullRequestCommentText(requester *Account, reviewers []*Account, pr *github.PullRequestEvent) string {
	users := make([]string, len(reviewers))
	for i, u := range reviewers {
		users[i] = u.GetUserPermalink()
	}

	return fmt.Sprintf(`<body>📋 <code>[<b>%s</b>] <a href="%s">Pull request #%d: %s</a> by %s

<b>%d</b> changed files (<b>+%d -%d</b>)
Reviewers: %s

by %s
</code></body>`,
		pr.PullRequest.GetState(),
		pr.PullRequest.GetHTMLURL(), pr.PullRequest.GetNumber(), pr.PullRequest.GetTitle(), requester.GetUserPermalink(),
		pr.PullRequest.GetChangedFiles(), pr.PullRequest.GetAdditions(), pr.PullRequest.GetDeletions(),
		strings.Join(users, ", "),
		signature,
	)
}

func createReviewRequestDescText(requester *Account, pr *github.PullRequestEvent) string {
	return fmt.Sprintf(`<body><a href="%s">#%d: %s</a> by %s

Could you please review a pull request 🙇

After you finished a code review, pass this assignee back to %s.
Do not mark complete.
</body>`,
		pr.PullRequest.GetHTMLURL(), pr.PullRequest.GetNumber(), pr.PullRequest.GetTitle(), requester.GetUserPermalink(),
		requester.GetUserPermalink(),
	)
}
