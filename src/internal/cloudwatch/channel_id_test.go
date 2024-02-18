package cloudwatch

import (
	"context"
	"fmt"
	"testing"

	"github.com/dstrates/cloudwatch-slack-alerts/src/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestUnitParseChannelMap(t *testing.T) {
	mockCfg := &config.Config{
		FeatureFlags: config.FeatureFlags{
			EnableSlackChannelMap: true,
		},
	}

	testCases := []struct {
		name              string
		functionName      string
		channelMapJSON    string
		expectedChannelID string
		expectedErr       string
	}{
		{
			name:              "Valid JSON input",
			functionName:      "function-1",
			channelMapJSON:    `{"function-1": "channel-1", "function-2": "channel-2"}`,
			expectedChannelID: "channel-1",
			expectedErr:       "",
		},
		{
			name:              "Valid JSON input",
			functionName:      "service-1-function-1",
			channelMapJSON:    `{"service-1": "channel-1", "service-1-function-1": "channel-2"}`,
			expectedChannelID: "channel-2",
			expectedErr:       "",
		},
		{
			name:              "Empty JSON input",
			functionName:      "function-1",
			channelMapJSON:    `{}`,
			expectedChannelID: "default",
			expectedErr:       "",
		},
		{
			name:              "Invalid JSON input",
			functionName:      "function-1",
			channelMapJSON:    `{"key1": "value1", "key2": "value2",}`,
			expectedChannelID: "",
			expectedErr:       "invalid character '}' looking for beginning of object key string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCfg.Slack.ChannelMap = tc.channelMapJSON
			mockCfg.Slack.DefaultChannelID = "default"
			result, err := cloudwatch.FindMatchingChannelID(context.TODO(), mockCfg, tc.functionName)
			assert.Equal(t, result, tc.expectedChannelID, fmt.Sprintf("test case: %s", tc.name))
			if tc.expectedErr == "" {
				assert.Nil(t, err, fmt.Sprintf("test case: %s", tc.name))
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr, fmt.Sprintf("test case: %s", tc.name))
			}
		})
	}
}
