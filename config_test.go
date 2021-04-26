package githubasana

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	data, err := os.ReadFile("./testdata/config.yml")
	if err != nil {
		t.Fatal(err)
	}

	c, err := LoadConfig(data)
	require.NoError(t, err)

	assert.Equal(t, "123", c.Accounts["user1"])
	assert.Equal(t, "456", c.Accounts["user2"])
	assert.Equal(t, "789", c.Accounts["user3"])
}
