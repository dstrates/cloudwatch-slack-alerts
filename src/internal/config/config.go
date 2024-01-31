package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/nullify-platform/logger/pkg/logger"
)

func GetDevelopmentMode() bool {
	return getStringVariable("DEVELOPMENT_MODE", "") == "true"
}

func GetLogLevel() string {
	return getStringVariable("LOG_LEVEL", "info")
}

func GetConfig(ctx context.Context) *Config {
	cfg := &Config{
		AWSRegion: getStringVariable("AWS_REGION", ""),
		Slack: Slack{
			DefaultChannelID: getStringVariable("DEFAULT_SLACK_CHANNEL_ID", ""),
		},
		FeatureFlags: FeatureFlags{
			EnableSlackChannelMap: getBooleanVariable("ENABLE_SLACK_CHANNEL_MAP", false),
		},
	}

	return cfg
}

type Config struct {
	AWSRegion string

	Slack        Slack
	FeatureFlags FeatureFlags
}

type Slack struct {
	DefaultChannelID string

	// lazy loaded fields
	ChannelMap string
	APIKey     string
}

type FeatureFlags struct {
	EnableSlackChannelMap bool
}

func (c *Config) Validate() error {
	if c.AWSRegion == "" {
		return errors.New("AWS_REGION is empty")
	}

	return nil
}

func (c *Config) GetSlackKey(ctx context.Context) (string, error) {
	if c.Slack.APIKey == "" {
		apiKey, err := getParameterStoreVariable(ctx, "SLACK_KEY_PARAMETER_NAME")
		if err != nil {
			return "", err
		}

		c.Slack.APIKey = apiKey
	}

	return c.Slack.APIKey, nil
}

func (c *Config) GetSlackChannelMapParameter(ctx context.Context) (string, error) {
	logger.Debug(
		"fetching slack channel map from SSM",
		logger.String("c.Slack.ChannelMap", c.Slack.ChannelMap),
	)
	fmt.Printf("fetching slack channel map from SSM: %s\n", c.Slack.ChannelMap)
	if c.Slack.ChannelMap == "" {
		channelMap, err := getParameterStoreVariable(ctx, "SLACK_CHANNEL_MAP_PARAMETER_NAME")
		if err != nil {
			return "", err
		}

		c.Slack.ChannelMap = channelMap
	}

	return c.Slack.ChannelMap, nil
}

func (c *Config) GetSlackDefaultChannelID(ctx context.Context) (string, error) {
	if c.Slack.DefaultChannelID == "" {
		defaultChannelID, err := getParameterStoreVariable(ctx, "DEFAULT_SLACK_CHANNEL_PARAMETER_NAME")
		if err != nil {
			return "", err
		}

		c.Slack.DefaultChannelID = defaultChannelID
	}

	return c.Slack.DefaultChannelID, nil
}
