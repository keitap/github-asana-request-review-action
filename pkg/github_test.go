package pkg

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadTestdata(filepath string) (name string, payload []byte) {
	s := strings.SplitN(path.Base(filepath), "-", 2)

	payload, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	return s[0], payload
}

func TestParseRequestReviewerEvent(t *testing.T) {
	name, payload := loadTestdata("testdata/pull_request-review_requested.json")
	e, err := parseRequestReviewerEvent(name, payload)
	require.NoError(t, err)

	assert.Equal(t, "keitap-2nd", *e.RequestedReviewer.Login)
}
