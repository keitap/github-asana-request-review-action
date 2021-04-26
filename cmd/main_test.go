package main

import (
	"strings"
	"testing"

	"github.com/google/go-github/v35/github"
	"github.com/stretchr/testify/assert"
)

func TestGetRepoFile(t *testing.T) {
	gh := github.NewClient(nil)
	data, err := getRepoFile(gh, "golang/go", "README.md", "6f3da9d2f6b4f7dbbe5d15260d87ed2a84488fde")
	if err != nil {
		t.Fatal(err)
	}

	str := string(*data)
	assert.True(t, strings.HasPrefix(str, "# The Go Programming Language"))
}
