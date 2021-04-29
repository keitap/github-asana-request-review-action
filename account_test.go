package githubasana

import (
	"testing"

	"bitbucket.org/mikehouston/asana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	asanaUserGID = "5590853215184"
)

func TestNewNoAsanaAccount(t *testing.T) {
	a := NewNoAsanaAccount("keitap")

	assert.Equal(t, "", a.AsanaUserGID)
	assert.Equal(t, "keitap", a.Name)
	assert.Equal(t, "keitap", a.GitHubLogin)
	assert.Equal(t, `<a href="https://github.com/keitap">keitap</a>`, a.GetUserPermalink())
}

func TestNewFromUserTaskListGID(t *testing.T) {
	c := asana.NewClientWithAccessToken(asanaToken)

	a, err := NewAccount(c, asanaUserGID, "keitap")
	require.NoError(t, err)

	assert.Equal(t, "5590853215184", a.AsanaUserGID)
	assert.Equal(t, "Keita Kitamura", a.Name)
	assert.Equal(t, "keitap", a.GitHubLogin)
	assert.Equal(t, `<a data-asana-gid="5590853215184"/>`, a.GetUserPermalink())
}
