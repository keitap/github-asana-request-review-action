package githubasana

import (
	"errors"
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

	log.Printf("asana: https://app.asana.com/0/%s/%s", projectID, taskID)

	requester := pr.PullRequest.User.GetLogin()
	reviewer := pr.RequestedReviewer.GetLogin()

	requesterAsanaGID := conf.Accounts[requester]
	if requesterAsanaGID == "" {
		return errors.New("requester asana GID is not set")
	}

	reviewerAsanaGID := conf.Accounts[reviewer]
	if reviewerAsanaGID == "" {
		return errors.New("reviewer asana GID is not set")
	}

	due := asana.Date(time.Now().AddDate(0, 0, 3))

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

	// add a review request task as a subtask if not exists.
	subtask, err := FindSubtaskByName(ac, taskID, reviewer)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	var subtaskStory *asana.Story

	if subtask == nil {
		subtask, err = AddCodeReviewSubtask(ac, taskID, requesterAsanaGID, reviewerAsanaGID, due, pr)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}

		subtaskStory = nil
	} else {
		subtaskStory, err = FindTaskComment(ac, subtask.ID, signature)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}

	// upsert a review description comment of a subtask.
	if err := upsertPullRequestComment(ac, subtask.ID, subtaskStory, pr); err != nil {
		return err
	}

	return nil
}

func upsertPullRequestComment(ac *asana.Client, taskID string, story *asana.Story, pr *github.PullRequestEvent) error {
	if story == nil {
		if _, err := AddPullRequestURLToTaskComment(ac, taskID, pr); err != nil {
			return xerrors.Errorf(": %w", err)
		}
	} else {
		if _, err := UpdateTaskComment(ac, story.ID, pr); err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}

	return nil
}
