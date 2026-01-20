package models

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMessageRequest_Validate_SwitchVersion(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		wantError bool
		errorType string
	}{
		{
			name:      "Valid text",
			text:      "Hello World",
			wantError: false,
		},
		{
			name:      "Empty text",
			text:      "",
			wantError: true,
			errorType: "empty",
		},
		{
			name:      "Too long text",
			text:      string(make([]byte, 5001)),
			wantError: true,
			errorType: "length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateMessageRequest{Text: tt.text}
			err := req.Validate()

			switch {
			case tt.wantError && tt.errorType == "empty":
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot be empty")
			case tt.wantError && tt.errorType == "length":
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "less than 5000")
			default:
				assert.NoError(t, err)
				assert.Equal(t, strings.TrimSpace(tt.text), req.Text)
			}
		})
	}
}
