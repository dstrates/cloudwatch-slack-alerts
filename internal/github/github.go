package github

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	DefaultOrg    = "Nullify-Platform"
	DefaultBranch = "main"
)

var repoFunctionMap = map[string]string{
	"force": "the-force",
	"sith":  "the-sith",
}

func ExtractPathAndLineNumber(caller string) (string, int, error) {
	parts := strings.Split(caller, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid caller format: %s", caller)
	}

	filePath := parts[0]
	lineStr := parts[1]

	lineNumber, err := strconv.Atoi(lineStr)
	if err != nil {
		return "", 0, fmt.Errorf("failed to convert line number: %s", err)
	}

	return filePath, lineNumber, nil
}

func FindMatchingRepo(repoFunctionMap map[string]string, functionName string) (string, bool) {
	functionName = strings.ToLower(functionName)
	parts := strings.Split(functionName, "-")

	for fn, repoName := range repoFunctionMap {
		lowerFn := strings.ToLower(fn)

		for _, part := range parts {
			if strings.Contains(lowerFn, part) {
				return repoName, true
			}
		}
	}

	return "", false
}

func validateParameters(org, branch, caller, functionName string) error {
	if org == "" || branch == "" || caller == "" || functionName == "" {
		return fmt.Errorf("missing or empty parameter(s)")
	}
	return nil
}

func ConstructGitHubURL(org, branch, caller, functionName string) (string, error) {
	if err := validateParameters(org, branch, caller, functionName); err != nil {
		return "", err
	}

	matchingRepo, found := FindMatchingRepo(repoFunctionMap, functionName)
	if !found {
		return "", fmt.Errorf("repo name not found for caller: %s. function: %s, repoFunctionMap: %+v", caller, functionName, repoFunctionMap)
	}

	urlFormat := "https://github.com/%s/%s/tree/%s/"
	gitHubURL := fmt.Sprintf(urlFormat, org, matchingRepo, branch)

	return gitHubURL, nil
}
