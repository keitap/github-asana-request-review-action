package pkg

import (
	"fmt"
	"regexp"
	"strings"

	"bitbucket.org/mikehouston/asana-go"
	"github.com/google/go-github/v35/github"
)

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
		HTMLText: fmt.Sprintf(`<body>Pull Request:
<a href="%s">%s</a>

by <b>%s</b>, <b>%d</b> changed files (<b>+%d -%d</b>)
</body>`,
			*pr.PullRequest.HTMLURL, *pr.PullRequest.HTMLURL,
			*pr.Sender.Login,
			*pr.PullRequest.ChangedFiles, *pr.PullRequest.Additions, *pr.PullRequest.Deletions,
		),
	}

	return task.CreateComment(client, story)
}

func HasCommentContainsURL(client *asana.Client, taskID string, url string) (bool, error) {
	task := &asana.Task{ID: taskID}
	stories, _, err := task.Stories(client, &asana.Options{
		Limit: 100,
	})
	if err != nil {
		return false, err
	}

	for _, s := range stories {
		if strings.Contains(s.Text, url) {
			return true, nil
		}
	}

	return false, nil
}
