package cloudwatch

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/dstrates/cloudwatch-slack-alerts/internal/config"
	"github.com/nullify-platform/logger/pkg/logger"
)

// FindMatchingChannelID finds the slack channel for the given function
// by longest prefix in the channel map or the default channel if no match is found
func FindMatchingChannelID(ctx context.Context, cfg *config.Config, functionName string) (string, error) {
	if cfg.FeatureFlags.EnableSlackChannelMap {
		channelMapValue, err := cfg.GetSlackChannelMapParameter(ctx)
		if err != nil {
			logger.Error("error fetching channel map SSM parameter", logger.Err(err))
			return "", err
		}

		channelID, err := findChannelIDInMap(channelMapValue, functionName)
		if err != nil {
			return "", err
		}

		// if a channel ID was found, return it
		// otherwise use the default one
		if channelID != "" {
			logger.Debug(
				"slack channel map found",
				logger.String("channelID", channelID),
			)

			return channelID, nil
		}
	}

	defaultChannelID, err := cfg.GetSlackDefaultChannelID(ctx)
	if err != nil {
		logger.Error("error fetching default slack channel from SSM", logger.Err(err))
		return "", err
	}

	logger.Debug(
		"default slack channel selected",
		logger.String("channelID", defaultChannelID),
	)

	return defaultChannelID, nil
}

func findChannelIDInMap(channelMapValue string, functionName string) (string, error) {
	channelMap := map[string]string{}

	err := json.Unmarshal([]byte(channelMapValue), &channelMap)
	if err != nil {
		return "", err
	}

	// find the channel ID with the longest prefix match on the function name
	channeID := ""
	longestPrefix := ""

	for prefix, channel := range channelMap {
		if strings.HasPrefix(functionName, prefix) && len(prefix) > len(longestPrefix) {
			channeID = channel
			longestPrefix = prefix
		}
	}

	return channeID, nil
}
