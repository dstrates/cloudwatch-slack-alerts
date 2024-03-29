package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractPathAndLineNumber(t *testing.T) {
	testCases := []struct {
		caller       string
		expectedPath string
		expectedLine int
		expectError  bool
	}{
		{
			caller:       "file.go:123",
			expectedPath: "file.go",
			expectedLine: 123,
			expectError:  false,
		},
		{
			caller:       "invalid_format",
			expectedPath: "",
			expectedLine: 0,
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		path, line, err := github.ExtractPathAndLineNumber(tc.caller)

		if tc.expectError {
			assert.Error(t, err, "Expected an error")
		} else {
			assert.NoError(t, err, "Expected no error")
			assert.Equal(t, tc.expectedPath, path, "Path should match")
			assert.Equal(t, tc.expectedLine, line, "Line number should match")
		}
	}
}

func TestFindPartialMatch(t *testing.T) {
	repoFunctionMap := map[string]string{
		"force": "the-force",
		"jedi":  "jedi-order",
	}

	testCases := []struct {
		partialRepoName string
		expectedRepo    string
		expectFound     bool
	}{
		{
			partialRepoName: "force",
			expectedRepo:    "the-force",
			expectFound:     true,
		},
		{
			partialRepoName: "jedi",
			expectedRepo:    "jedi-order",
			expectFound:     true,
		},
		{
			partialRepoName: "unknown",
			expectedRepo:    "",
			expectFound:     false,
		},
	}

	for _, tc := range testCases {
		repo, found := github.FindMatchingRepo(repoFunctionMap, tc.partialRepoName)

		assert.Equal(t, tc.expectedRepo, repo, "Repo should match")
		assert.Equal(t, tc.expectFound, found, "Found should match")
	}
}

func TestConstructGitHubURL(t *testing.T) {
	testCases := []struct {
		org         string
		branch      string
		caller      string
		function    string
		expectedURL string
		expectError bool
	}{
		{
			org:         "dstrates",
			branch:      "main",
			caller:      "cmd/alerter/main.go:59",
			function:    "cloudwatch-slack-alerts",
			expectedURL: "https://github.com/dstrates/cloudwatch-slack-alerts/tree/main/cmd/alerter/main.go",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		url, err := github.ConstructGitHubURL(tc.org, tc.branch, tc.caller, tc.function)

		if tc.expectError {
			assert.Error(t, err, "Expected an error")
		} else {
			assert.NoError(t, err, "Expected no errors")
			assert.Equal(t, tc.expectedURL, url, "URL should match")
			t.Logf("Constructed URL: %s", url)
		}
	}
}
