package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dstrates/cloudwatch-slack-alerts/src/internal/cloudwatch"
	"github.com/dstrates/cloudwatch-slack-alerts/src/internal/config"
	"github.com/dstrates/cloudwatch-slack-alerts/src/internal/github"
	"github.com/dstrates/cloudwatch-slack-alerts/src/internal/slack"
	"github.com/nullify-platform/logger/pkg/logger"
)

func main() {
	lambda.Start(handler)
}
func handler(ctx context.Context, event events.CloudwatchLogsEvent) error {
	var log logger.Logger
	var err error

	if config.GetDevelopmentMode() {
		log, err = logger.ConfigureDevelopmentLogger(config.GetLogLevel())
	} else {
		log, err = logger.ConfigureProductionLogger(config.GetLogLevel())
	}

	if err != nil {
		logger.Error("error configuring logger", logger.Err(err))
		return err
	}
	defer log.Sync()

	cfg := config.GetConfig(ctx)

	if err := cfg.Validate(); err != nil {
		logger.Error("invalid config", logger.Err(err))
		return err
	}

	logDetails, logMessage, err := cloudwatch.ProcessCloudWatchLog(ctx, event, cfg)
	if err != nil {
		logger.Error("error processing cloudwatch log", logger.Err(err))
		return err
	}

	logger.Debug(
		"cloudwatch log successfully processed",
		logger.Any("logDetails", logDetails),
		logger.Any("logMessage", logMessage),
	)

	channelID, err := cloudwatch.FindMatchingChannelID(ctx, cfg, logDetails.FunctionName)
	if err != nil {
		logger.Error("error finding matching Slack channel ID", logger.Err(err))
		return err
	}

	gitHubURL, err := github.ConstructGitHubURL(github.DefaultOrg, github.DefaultBranch, logMessage.Caller, logDetails.FunctionName)
	if err != nil {
		logger.Error("error constructing GitHub URL", logger.Err(err))
	}

	err = slack.PostToSlack(ctx, cfg, logDetails, logMessage, channelID, gitHubURL)
	if err != nil {
		logger.Error("error posting message to Slack channel", logger.Err(err))
		return err
	}

	return nil
}
