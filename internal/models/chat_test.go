package models

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateChatRequest_Validate_SwitchVersion(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		wantError bool
		errorType string
	}{
		{
			name:      "Valid title",
			title:     "Valid Chat",
			wantError: false,
		},
		{
			name:      "Empty title",
			title:     "",
			wantError: true,
			errorType: "empty",
		},
		{
			name:      "Title with only spaces",
			title:     "   ",
			wantError: true,
			errorType: "empty",
		},
		{
			name:      "Too long title",
			title:     string(make([]byte, 201)),
			wantError: true,
			errorType: "length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateChatRequest{Title: tt.title}
			err := req.Validate()

			switch {
			case tt.wantError && tt.errorType == "empty":
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot be empty")
			case tt.wantError && tt.errorType == "length":
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "less than 200")
			default:
				assert.NoError(t, err)
				assert.Equal(t, strings.TrimSpace(tt.title), req.Title)
			}
		})
	}
}
