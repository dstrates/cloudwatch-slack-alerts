package cloudwatch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dstrates/cloudwatch-slack-alerts/src/internal/config"
	"github.com/nullify-platform/logger/pkg/logger"
)

type LogDetails struct {
	Owner         string `json:"owner"`
	Region        string `json:"region"`
	ErrorType     string `json:"errorType"`
	FunctionName  string
	LogGroupName  string `json:"logGroup"`
	LogStreamName string `json:"logStream"`
	LogStreamURL  string
}

type LogMessage struct {
	Level             string  `json:"level"`
	Timestamp         float64 `json:"ts"`
	Caller            string  `json:"caller"`
	Msg               string  `json:"msg"`
	Version           string  `json:"version"`
	TenantID          string  `json:"tenantId"`
	Platform          string  `json:"platform"`
	Repo              string  `json:"repo"`
	Repository        string  `json:"repository"`
	PullRequestNumber int     `json:"pullRequestNumber"`
	CloneURL          string  `json:"cloneUrl"`
	SourceBranch      string  `json:"sourceBranch"`
	SourceCommit      string  `json:"sourceCommit"`
	TargetBranch      string  `json:"targetBranch"`
	TargetCommit      string  `json:"targetCommit"`
	NumTools          int     `json:"numTools"`
	CheckID           int     `json:"checkId"`
	Language          string  `json:"language"`
	Tools             string  `json:"tools"`
	Scanner           string  `json:"scanner"`
	Error             string  `json:"error"`
	ErrorVerbose      string  `json:"errorVerbose"`
	StdError          string  `json:"stdError"`
}

func ProcessCloudWatchLog(ctx context.Context, event events.CloudwatchLogsEvent, cfg *config.Config) (*LogDetails, *LogMessage, error) {
	logsData, err := parseCloudwatchLogsData(event)
	if err != nil {
		logger.Error("error parsing CloudWatch Logs data", logger.Err(err))
		return nil, nil, err
	}

	logMessage, err := parseLogMessage(logsData)
	if err != nil {
		logger.Error("error parsing log message", logger.Err(err))
		return nil, nil, err
	}

	errorType := parseErrorType(logsData)

	functionName, err := parseFunctionName(logsData.LogGroup)
	if err != nil {
		logger.Error("error parsing function name", logger.Err(err))
		return nil, nil, err
	}

	logStreamURL := constructLogStreamURL(cfg.AWSRegion, logsData.LogGroup, logsData.LogStream)

	logDetails := LogDetails{
		Owner:         logsData.Owner,
		Region:        cfg.AWSRegion,
		ErrorType:     errorType,
		FunctionName:  functionName,
		LogGroupName:  logsData.LogGroup,
		LogStreamName: logsData.LogStream,
		LogStreamURL:  logStreamURL,
	}

	return &logDetails, logMessage, nil
}

// parseCloudwatchLogsData parses the base64 encoded raw data into CloudwatchLogsData
func parseCloudwatchLogsData(event events.CloudwatchLogsEvent) (*events.CloudwatchLogsData, error) {
	logger.Debug(
		"raw cloudwatch log data",
		logger.Any("rawData", event.AWSLogs),
	)

	parsedData, err := event.AWSLogs.Parse()
	if err != nil {
		logger.Error(
			"error parsing raw cloudwatch log data",
			logger.Err(err),
		)
		return nil, fmt.Errorf("error parsing cloudwatch data: %w", err)
	}

	logger.Debug(
		"cloudwatch log data parsed succesfully",
		logger.Any("parsedData", parsedData),
	)

	return &parsedData, nil
}

// parseLogMessage parses the JSON log event into a LogMessage struct
func parseLogMessage(parsedData *events.CloudwatchLogsData) (*LogMessage, error) {
	if len(parsedData.LogEvents) == 0 {
		return nil, errors.New("no log events found in parsedData")
	}

	message := parsedData.LogEvents[0].Message

	var logMessage LogMessage
	if err := json.Unmarshal([]byte(message), &logMessage); err != nil {
		return nil, err
	}

	return &logMessage, nil
}

// parseFunctionName parses the CloudWatch log group name for the lambda function name e.g. /aws/lambda/my-function-name
func parseFunctionName(logGroup string) (string, error) {
	parts := strings.Split(logGroup, "/")
	if len(parts) == 4 {
		functionName := parts[3]
		logger.Info("function name parsed succesfully",
			logger.String("FunctionName", functionName),
		)
		return functionName, nil
	}
	return "", fmt.Errorf("unexpected number of parts: %v", parts)
}

// parseErrorType parses a log event and returns the error type as a string
func parseErrorType(logsData *events.CloudwatchLogsData) string {
	for _, logEvent := range logsData.LogEvents {
		message := logEvent.Message
		switch {
		case strings.Contains(message, "Task timed out"):
			return ":alarm_clock: Timeout Error"
		case strings.Contains(message, "\"level\":\"error\""):
			return ":exclamation: Application Error"
		case strings.Contains(message, "Invocation error"),
			strings.Contains(message, "AccessDeniedException"),
			strings.Contains(message, "ResourceNotFoundException"):
			return ":lock: Invocation Error"
		case strings.Contains(message, "panic: runtime error:"),
			strings.Contains(message, "OutOfMemoryError"),
			strings.Contains(message, "NetworkError"):
			return ":boom: Runtime Error"
		}
	}

	return ":grey_question: Unknown Error"
}

// constructLogStreamURL creates the URL to the log stream event for Slack
func constructLogStreamURL(awsRegion, logGroup, logStream string) string {
	// Replace '/' with '$252F' in the log group
	logGroup = strings.ReplaceAll(logGroup, "/", "$252F")

	return fmt.Sprintf(
		"https://console.aws.amazon.com/cloudwatch/home?region=%s#logsV2:log-groups/log-group/%s/log-events/%s",
		awsRegion, logGroup, url.PathEscape(logStream),
	)
}
