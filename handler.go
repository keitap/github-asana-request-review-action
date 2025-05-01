package githubasana

import (
	"log"
	"time"

	"bitbucket.org/mikehouston/asana-go"
	"github.com/google/go-github/v71/github"
	"golang.org/x/xerrors"
)

const (
	prEventActionSynchronize            = "synchronize"
	prEventActionEdited                 = "edited"
	prEventActionLabeled                = "labeled"
	prEventActionUnlabeled              = "unlabeled"
	prEventActionReviewRequested        = "review_requested"
	prEventActionReviewRequestedRemoved = "review_request_removed"
)

type Handler struct {
	conf *Config
	ac   *asana.Client
	gh   *github.Client
}

func NewHandler(conf *Config, asanaClient *asana.Client, githubClient *github.Client) *Handler {
	return &Handler{
		conf: conf,
		ac:   asanaClient,
		gh:   githubClient,
	}
}

func (h *Handler) Handle(eventName string, eventPayload []byte) error {
	event, err := github.ParseWebHook(eventName, eventPayload)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	switch e := event.(type) {
	case *github.PullRequestEvent:
		return h.handlePullRequestEvent(e)
	case *github.PullRequestReviewEvent:
		return h.handlePullRequestReviewEvent(e)
	default:
		log.Println("unknown event: " + eventName)
	}

	return nil
}

func (h *Handler) handlePullRequestEvent(pr *github.PullRequestEvent) error {
	workspaceID, projectID, taskID := parseAsanaTaskLink(pr.PullRequest.GetBody())
	if taskID == "" {
		log.Println("asana task url not found in description.")

		return nil
	}

	if projectID != "" {
		log.Printf("asana: https://app.asana.com/0/%s/%s", projectID, taskID)
	} else {
		log.Printf("asana: https://app.asana.com/1/%s/task/%s?focus=true", workspaceID, taskID)
	}

	requester, err := h.fetchAccount(pr.PullRequest.User.GetLogin())
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	log.Printf("requester: %s", requester)

	// add a review description comment to a parent task if not exists.
	if err := h.updateTask(pr, requester, taskID); err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return h.updateReviewers(pr, requester, taskID)
}

func (h *Handler) updateTask(pr *github.PullRequestEvent, requester *Account, taskID string) error {
	// add a review description comment to a parent task if not exists.
	story, err := FindTaskComment(h.ac, taskID, signature)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	// upsert a review description comment of a parent task.
	if err := h.upsertPullRequestComment(taskID, story, requester, pr); err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}

func (h *Handler) updateReviewers(pr *github.PullRequestEvent, requester *Account, taskID string) error {
	var ghReviewers []*github.User
	var reviewers []*Account

	hasRequestedReviewersFields := pr.GetAction() == prEventActionReviewRequested ||
		pr.GetAction() == prEventActionReviewRequestedRemoved

	shouldUpdateReviewerEvent := hasRequestedReviewersFields ||
		pr.GetAction() == prEventActionLabeled ||
		pr.GetAction() == prEventActionUnlabeled ||
		// to update diff numbers.
		pr.GetAction() == prEventActionSynchronize ||
		// in case when a reviewer is assigned before an Asana URL is added to the pull request description.
		pr.GetAction() == prEventActionEdited

	if !shouldUpdateReviewerEvent {
		return nil
	}

	if hasRequestedReviewersFields {
		ghReviewers = pr.PullRequest.RequestedReviewers
	} else {
		var err error
		ghReviewers, err = getRequestedReviewers(h.gh, pr.GetRepo().GetOwner().GetLogin(), pr.GetRepo().GetName(), pr.GetNumber())
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}

	reviewers = make([]*Account, len(ghReviewers))

	for i, r := range ghReviewers {
		var err error
		reviewers[i], err = h.fetchAccount(r.GetLogin())
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}

		log.Printf("reviewer: %s", reviewers[i])
	}

	for _, reviewer := range reviewers {
		err := h.upsertReviewer(pr, requester, reviewer, taskID)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}

	return nil
}

func (h *Handler) upsertReviewer(pr *github.PullRequestEvent, requester *Account, reviewer *Account, taskID string) error {
	if reviewer.AsanaUserGID == "" {
		log.Printf("reviewer has no asana account: %s", reviewer.GitHubLogin)

		return nil
	}

	due := asana.Date(NextBusinessDay(h.conf.DueDate, time.Now(), h.conf.Holidays))

	// add a review request task as a subtask if not exists.
	subtask, err := FindSubtaskByName(h.ac, taskID, reviewer.Name)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if subtask == nil {
		log.Printf("code review subtask not found. will create one.")

		subtask, err = AddCodeReviewSubtask(h.ac, taskID, pr.GetNumber(), requester, reviewer, due, pr)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}

		log.Printf("added code review subtask to feature task: %s", subtask.ID)
	} else {
		err := UpdateCodeReviewSubtask(h.ac, subtask, requester, pr)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}

		log.Printf("update code review subtask: %s", subtask.ID)
	}

	return nil
}

func (h *Handler) upsertPullRequestComment(taskID string, story *asana.Story, requester *Account, pr *github.PullRequestEvent) error {
	if story == nil {
		if _, err := AddPullRequestCommentToTask(h.ac, taskID, requester, pr); err != nil {
			return xerrors.Errorf(": %w", err)
		}

		log.Printf("added comment to task: %s", taskID)
	} else {
		if _, err := UpdateTaskComment(h.ac, story.ID, requester, pr); err != nil {
			return xerrors.Errorf(": %w", err)
		}

		log.Printf("updated comment on task: %s %s", taskID, story.ID)
	}

	return nil
}

func (h *Handler) fetchAccount(githubLogin string) (*Account, error) {
	userGID := h.conf.Accounts[githubLogin]

	if userGID == "" {
		return NewNoAsanaAccount(githubLogin), nil
	}

	a, err := NewAccount(h.ac, userGID, githubLogin)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (h *Handler) handlePullRequestReviewEvent(pr *github.PullRequestReviewEvent) error {
	_, _, taskID := parseAsanaTaskLink(pr.PullRequest.GetBody())
	if taskID == "" {
		log.Println("asana task url not found in description.")

		return nil
	}

	requester, err := h.fetchAccount(pr.PullRequest.User.GetLogin())
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	reviewer, err := h.fetchAccount(pr.Review.User.GetLogin())
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	subtask, err := FindSubtaskByName(h.ac, taskID, reviewer.Name)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if subtask == nil {
		log.Printf("review subtask is not found.")

		return nil
	}

	_, err = AddCodeReviewSubtaskComment(h.ac, subtask, requester, reviewer, pr)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	log.Printf("review is submitted by reviewer: %s state: %s", reviewer.Name, pr.Review.GetState())

	return nil
}
