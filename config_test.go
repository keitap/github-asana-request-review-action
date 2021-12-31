package githubasana

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Run("config_v1.0.yml", func(t *testing.T) {
		data, err := os.ReadFile("./testdata/config_v1.0.yml")
		if err != nil {
			t.Fatal(err)
		}

		c, err := LoadConfig(data)
		require.NoError(t, err)

		assert.Equal(t, 1, c.DueDate)
		assert.Equal(t, "123", c.Accounts["user1"])
		assert.Equal(t, "456", c.Accounts["user2"])
		assert.Equal(t, "789", c.Accounts["user3"])
	})

	t.Run("config_v1.1.yml", func(t *testing.T) {
		data, err := os.ReadFile("./testdata/config_v1.1.yml")
		if err != nil {
			t.Fatal(err)
		}

		c, err := LoadConfig(data)
		require.NoError(t, err)

		assert.Equal(t, 3, c.DueDate)
		assert.Equal(t, true, c.Holidays["2021-12-31"])
		assert.Equal(t, false, c.Holidays["2022-01-01"])
		assert.Equal(t, false, c.Holidays["2022-01-02"])
	})
}
