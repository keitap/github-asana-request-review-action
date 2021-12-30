package githubasana

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNextBusinessDay(t *testing.T) {
	holidays := map[string]bool{
		"2021-05-01": false,
		"2021-05-03": true,
	}

	tests := []struct {
		date     string
		n        int
		expected string
	}{
		{"2021-05-06", 0, "2021-05-06"},
		{"2021-05-06", 1, "2021-05-07"},
		{"2021-05-06", 2, "2021-05-10"},
		{"2021-05-06", 3, "2021-05-11"},
		{"2021-05-06", 10, "2021-05-20"},
		{"2021-05-06", 100, "2021-09-23"},
		{"2021-05-06", 1000, "2025-03-06"},
		{"2021-05-08", 0, "2021-05-10"},
		{"2021-05-08", 1, "2021-05-10"},
		{"2021-05-08", 2, "2021-05-11"},
		{"2021-05-09", 0, "2021-05-10"},
		{"2021-05-09", 1, "2021-05-10"},
		{"2021-05-09", 2, "2021-05-11"},
		{"2021-05-09", 3, "2021-05-12"},
		//
		{"2021-05-01", 0, "2021-05-01"},
		{"2021-05-01", 1, "2021-05-04"},
		{"2021-05-01", 2, "2021-05-05"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s-%d", test.date, test.n), func(t *testing.T) {
			b, _ := time.Parse("2006-01-02", test.date)
			assert.Equal(t, test.expected, NextBusinessDay(test.n, b, holidays).Format("2006-01-02"))
		})
	}
}
