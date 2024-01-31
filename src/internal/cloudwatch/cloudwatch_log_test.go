package tests

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dstrates/cloudwatch-slack-alerts/internal/cloudwatch"
	"github.com/dstrates/cloudwatch-slack-alerts/internal/config"
	"github.com/stretchr/testify/require"
)

type LogEventMessage struct {
	Level           string  `json:"level"`
	Timestamp       float64 `json:"ts"`
	Caller          string  `json:"caller"`
	Msg             string  `json:"msg"`
	Version         string  `json:"version"`
	Platform        string  `json:"platform"`
	Repo            string  `json:"repo"`
	CloneURL        string  `json:"cloneUrl"`
	Branch          string  `json:"branch"`
	IsDefaultBranch bool    `json:"isDefaultBranch"`
	Commit          string  `json:"commit"`
	AuthorName      string  `json:"authorName"`
	Error           string  `json:"error"`
	ErrorVerbose    string  `json:"errorVerbose"`
}

func TestParsingCloudwatchLog(t *testing.T) {
	ctx := context.TODO()

	logEventMessage := LogEventMessage{
		Level:           "error",
		Msg:             "publishing a test error message for cloudwatch alerts",
		Platform:        "GitHub",
		Repo:            "test/leaky-secrets",
		Branch:          "main",
		IsDefaultBranch: true,
		Commit:          "af7f90fc2da62d0ca9a4d89f00d88e2b2ff7b643",
		AuthorName:      "test",
		Error:           "test error message",
		ErrorVerbose:    "test error message-verbose",
	}

	logEventMessageData, err := json.Marshal(&logEventMessage)
	require.NoError(t, err)

	logData := events.CloudwatchLogsData{
		Owner:               "012345678901", // AWS Account ID
		LogGroup:            "/aws/lambda/test-function",
		LogStream:           "2023/09/29/[$LATEST]718af798bded4b048af83eb06e98d572",
		SubscriptionFilters: []string{"error-alerts"},
		LogEvents: []events.CloudwatchLogsLogEvent{
			{
				ID:        "37821166091237533543887056780244332882191868685233160195",
				Timestamp: 1695959742804,
				Message:   string(logEventMessageData),
			},
		},
	}

	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	err = json.NewEncoder(writer).Encode(&logData)
	require.NoError(t, err)
	err = writer.Close()
	require.NoError(t, err)

	logEvent := events.CloudwatchLogsEvent{
		AWSLogs: events.CloudwatchLogsRawData{
			Data: base64.StdEncoding.EncodeToString(buf.Bytes()),
		},
	}

	cfg := &config.Config{
		AWSRegion: "us-east-2",
	}

	logDetails, logMessage, err := cloudwatch.ProcessCloudWatchLog(ctx, logEvent, cfg)
	require.NoError(t, err)

	t.Logf("logDetails: %+v", logDetails)
	t.Logf("logMessage: %+v", logMessage)
}
