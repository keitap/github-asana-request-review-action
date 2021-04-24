package main

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

func main() {
	gh := getGithubClient(os.Getenv("GITHUB_TOKEN"))
	eventName := os.Getenv("GITHUB_EVENT_NAME")
	eventPayload := getEventPayload(os.Getenv("GITHUB_EVENT_PATH"))

	log.Println(gh)
	log.Println(eventName)
	log.Println(eventPayload)
}

func getGithubClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	client := oauth2.NewClient(ctx, ts)
	return github.NewClient(client)
}

func getEventPayload(path string) *[]byte {
	log.Printf("GITHUB_EVENT_PATH: %s", path)

	payload, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to load event payload from %s: %s", path, err)
	}

	return &payload
}
