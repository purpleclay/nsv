package ci_test

import (
	"testing"

	"github.com/purpleclay/nsv/internal/ci"
	"github.com/stretchr/testify/require"
)

func TestDetectSkipPipelineTag(t *testing.T) {
	tests := []struct {
		name     string
		env      []string
		expected string
	}{
		{
			name:     "Default",
			env:      []string{},
			expected: "[skip ci]",
		},
		{
			name:     "Drone",
			env:      []string{"DRONE", "true"},
			expected: "[CI SKIP]",
		},
		{
			name:     "Jenkins",
			env:      []string{"JENKINS_URL", "http://jenkins"},
			expected: "[ci skip]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < len(tt.env); i += 2 {
				t.Setenv(tt.env[i], tt.env[i+1])
			}

			actual := ci.Detect()
			require.Equal(t, tt.expected, actual.SkipPipelineTag)
		})
	}
}
