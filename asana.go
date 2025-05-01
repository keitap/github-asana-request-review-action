package githubasana

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"bitbucket.org/mikehouston/asana-go"
	"github.com/google/go-github/v71/github"
)

const signature = "#github-asana-request-review"

var taskURLMatcher = regexp.MustCompile(`https://app.asana.com/0/(\d+)/(\d+)`)

func parseAsanaTaskLink(text string) (workspaceID string, projectID string, taskID string) {
	v1 := regexp.MustCompile(`https://app\.asana\.com/1/(\d+)/(?:task/|project/\d+/task/)(\d+)`)
	m := v1.FindStringSubmatch(text)

	if 0 < len(m) {
		return m[1], "", m[2]
	}

	m = taskURLMatcher.FindStringSubmatch(text)
	if 0 < len(m) {
		return "", m[1], m[2]
	}

	return "", "", ""
}

func AddPullRequestCommentToTask(client *asana.Client, taskID string, requester *Account, pr *github.PullRequestEvent) (*asana.Story, error) {
	task := &asana.Task{ID: taskID}
	story := &asana.StoryBase{
		HTMLText: createPullRequestCommentText(requester, pr),
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

func UpdateTaskComment(client *asana.Client, storyID string, requester *Account, pr *github.PullRequestEvent) (*asana.Story, error) {
	story := &asana.Story{ID: storyID}
	newStory := &asana.StoryBase{
		HTMLText: createPullRequestCommentText(requester, pr),
		IsPinned: true,
	}

	return story.UpdateStory(client, newStory)
}

func AddCodeReviewSubtask(client *asana.Client, taskID string, prID int, requester *Account, reviewer *Account, dueDate asana.Date, pr *github.PullRequestEvent) (*asana.Task, error) {
	req := &asana.CreateTaskRequest{
		Parent:    taskID,
		Assignee:  reviewer.AsanaUserGID,
		Followers: []string{requester.AsanaUserGID},
		TaskBase: asana.TaskBase{
			Name:      fmt.Sprintf(`✍️ Code review: #%d %s`, prID, reviewer),
			HTMLNotes: createReviewRequestDescText(requester, pr),
			DueOn:     &dueDate,
		},
	}

	return client.CreateTask(req)
}

func UpdateCodeReviewSubtask(client *asana.Client, subtask *asana.Task, requester *Account, pr *github.PullRequestEvent) error {
	update := &asana.UpdateTaskRequest{
		TaskBase: asana.TaskBase{
			HTMLNotes: createReviewRequestDescText(requester, pr),
		},
	}

	return subtask.Update(client, update)
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

func AddCodeReviewSubtaskComment(client *asana.Client, subtask *asana.Task, requester *Account, reviewer *Account, pr *github.PullRequestReviewEvent) (*asana.Story, error) {
	state := ""
	switch pr.Review.GetState() {
	case "approved":
		state = "✅ Approved"
	case "commented":
		state = "💬 Commented"
	case "changes_requested":
		state = "❗️ Changes Requested"
	}

	story, err := subtask.CreateComment(client, &asana.StoryBase{
		HTMLText: fmt.Sprintf(`<body><a href="%s"><b>%s</b></a>: reviewed by %s</body>`,
			pr.Review.GetHTMLURL(), state, reviewer.GetUserPermalink()),
	})
	if err != nil {
		return nil, err
	}

	// reassign to requester.
	err = subtask.Update(client, &asana.UpdateTaskRequest{
		Assignee: requester.AsanaUserGID,
	})
	if err != nil {
		return nil, err
	}

	return story, nil
}

func createPullRequestCommentText(requester *Account, pr *github.PullRequestEvent) string {
	return fmt.Sprintf(`<body>📋 <code>[<b>%s</b>] <a href="%s">Pull request #%d: %s</a> by %s

<b>%d</b> changed files (<b>+%d -%d</b>)
%s

%s
</code></body>`,
		pr.PullRequest.GetState(),
		pr.PullRequest.GetHTMLURL(), pr.PullRequest.GetNumber(), pr.PullRequest.GetTitle(), requester.GetUserPermalink(),
		pr.PullRequest.GetChangedFiles(), pr.PullRequest.GetAdditions(), pr.PullRequest.GetDeletions(),
		getLabelsText(pr),
		signature,
	)
}

func createReviewRequestDescText(requester *Account, pr *github.PullRequestEvent) string {
	return fmt.Sprintf(`<body><a href="%s">#%d: %s</a> by %s

<b>%d</b> changed files (<b>+%d -%d</b>)
%s

Could you please review a pull request ❤️
</body>`,
		pr.PullRequest.GetHTMLURL(), pr.PullRequest.GetNumber(), pr.PullRequest.GetTitle(), requester.GetUserPermalink(),
		pr.PullRequest.GetChangedFiles(), pr.PullRequest.GetAdditions(), pr.PullRequest.GetDeletions(),
		getLabelsText(pr),
	)
}

func getLabelsText(pr *github.PullRequestEvent) string {
	labels := make([]string, 0)
	for _, l := range pr.PullRequest.Labels {
		labels = append(labels, fmt.Sprintf("#<b>%s</b>", l.GetName()))
	}

	if len(labels) == 0 {
		return ""
	}

	sort.Strings(labels)

	return "Labels: " + strings.Join(labels, ", ")
}
