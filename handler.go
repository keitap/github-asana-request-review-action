package githubasana

import (
	"fmt"
	"log"
	"os"
	"time"

	"bitbucket.org/mikehouston/asana-go"
	"github.com/google/go-github/v35/github"
	"golang.org/x/xerrors"
)

func Handle(conf *Config, eventName string, eventPayload []byte) error {
	event, err := github.ParseWebHook(eventName, eventPayload)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	switch e := event.(type) {
	case *github.PullRequestEvent:
		return handlePullRequestEvent(conf, e)
	default:
		log.Println("unknown event: " + eventName)
	}

	return nil
}

func handlePullRequestEvent(conf *Config, pr *github.PullRequestEvent) error {
	projectID, taskID := parseAsanaTaskLink(pr.PullRequest.GetBody())
	if projectID == "" || taskID == "" {
		log.Println("asana task url not found in description.")
		return nil
	}

	requester := pr.PullRequest.User.GetLogin()
	reviewer := pr.RequestedReviewer.GetLogin()

	log.Printf("asana: https://app.asana.com/0/%s/%s", projectID, taskID)
	log.Printf("requester: %s", requester)
	log.Printf("reviewer: %s", reviewer)

	ac := asana.NewClientWithAccessToken(os.Getenv("ASANA_TOKEN"))

	// add a review description comment to a parent task if not exists.
	story, err := FindTaskComment(ac, taskID, signature)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	// upsert a review description comment of a parent task.
	if err := upsertPullRequestComment(ac, taskID, story, pr); err != nil {
		return err
	}

	//
	if reviewer == "" {
		log.Println("reviewer is not specified.")
		return nil
	}

	requesterAsanaGID := conf.Accounts[requester]
	if requesterAsanaGID == "" {
		return fmt.Errorf("requester asana GID is not set: %s", requester)
	}

	reviewerAsanaGID := conf.Accounts[reviewer]
	if reviewerAsanaGID == "" {
		return fmt.Errorf("reviewer asana GID is not set: %s", reviewer)
	}

	due := asana.Date(time.Now().AddDate(0, 0, 3))

	// add a review request task as a subtask if not exists.
	subtask, err := FindSubtaskByName(ac, taskID, reviewer)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if subtask == nil {
		log.Printf("code review subtask not found. will create one.")

		subtask, err = AddCodeReviewSubtask(ac, taskID, requesterAsanaGID, reviewerAsanaGID, due, pr)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}

		log.Printf("added code review subtask to feature task: %s", subtask.ID)
	}

	return nil
}

func upsertPullRequestComment(ac *asana.Client, taskID string, story *asana.Story, pr *github.PullRequestEvent) error {
	if story == nil {
		if _, err := AddPullRequestCommentToTask(ac, taskID, pr); err != nil {
			return xerrors.Errorf(": %w", err)
		}

		log.Printf("added comment to task: %s", taskID)
	} else {
		if _, err := UpdateTaskComment(ac, story.ID, pr); err != nil {
			return xerrors.Errorf(": %w", err)
		}

		log.Printf("updated comment on task: %s %s", taskID, story.ID)
	}

	return nil
}
