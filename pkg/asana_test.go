package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAsanaTaskLink(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		projectID string
		taskID    string
	}{
		{
			name:      `full screen URL`,
			text:      `Here is the Asana task URL that you should read before doing code review.\r\nhttps://app.asana.com/0/364167036366785/1162650948650897/f\r\n`,
			projectID: `364167036366785`,
			taskID:    `1162650948650897`,
		},
		{
			name:      `no full screen URL`,
			text:      `Task URL: https://app.asana.com/0/364167036366785/1162650948650897`,
			projectID: `364167036366785`,
			taskID:    `1162650948650897`,
		},
		{
			name:      `No task URL`,
			text:      `Lorem ipsum dolor sit amet, consectetur adipiscing elit.`,
			projectID: ``,
			taskID:    ``,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			projectID, taskID := parseAsanaTaskLink(test.text)
			assert.Equal(t, test.projectID, projectID)
			assert.Equal(t, test.taskID, taskID)
		})
	}
}
