package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v35/github"
	githubasana "github.com/keitap/github-asana-request-review-action"
	"golang.org/x/oauth2"
)

func main() {
	gh := getGithubClient(os.Getenv("GITHUB_TOKEN"))
	eventName := os.Getenv("GITHUB_EVENT_NAME")
	eventPayload := getEventPayload(os.Getenv("GITHUB_EVENT_PATH"))

	log.Printf("github event: %s", eventName)

	configPath := os.Getenv("INPUT_CONFIG_PATH")

	log.Printf("config path: %s", configPath)

	configData, err := getRepoFile(
		gh,
		os.Getenv("GITHUB_REPOSITORY"),
		configPath,
		os.Getenv("GITHUB_SHA"))
	if err != nil {
		log.Printf("::error ::cannot get config file: %s", err)
		os.Exit(1)
	}

	conf, err := githubasana.LoadConfig(*configData)
	if err != nil {
		log.Printf("::error ::cannot load config file: %s", err)
		os.Exit(1)
	}

	if err := githubasana.Handle(conf, eventName, *eventPayload); err != nil {
		log.Printf("::warning ::cannot handle github event: %s", err)
	}

	log.Printf("done.")
}

func getGithubClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	client := oauth2.NewClient(ctx, ts)
	return github.NewClient(client)
}

func getEventPayload(path string) *[]byte {
	payload, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to load event payload from %s: %s", path, err)
	}

	return &payload
}

func getRepoFile(gh *github.Client, repo, file, sha string) (*[]byte, error) {
	t := strings.Split(repo, "/")
	owner, repoName := t[0], t[1]

	fileContent, _, _, err := gh.Repositories.GetContents(
		context.Background(),
		owner,
		repoName,
		file,
		&github.RepositoryContentGetOptions{Ref: sha})

	var content string
	if err == nil {
		content, err = fileContent.GetContent()
	}

	if err != nil {
		log.Printf("Unable to load file from %s@%s/%s: %s", repo, sha, file, err)
		return nil, err
	}

	raw := []byte(content)
	return &raw, err
}
