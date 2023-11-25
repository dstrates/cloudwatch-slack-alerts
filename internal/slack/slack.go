package slack

import (
	"context"
	"fmt"

	"github.com/dstrates/cloudwatch-slack-alerts/internal/cloudwatch"
	"github.com/dstrates/cloudwatch-slack-alerts/internal/config"
	"github.com/nullify-platform/logger/pkg/logger"
	"github.com/slack-go/slack"
)

func PostToSlack(ctx context.Context, cfg *config.Config, logDetails *cloudwatch.LogDetails, logMessage *cloudwatch.LogMessage, channelID, githubURL string) error {
	apiKey, err := cfg.GetSlackKey(ctx)
	if err != nil {
		logger.Error("error fetching Slack API key", logger.Err(err))
		return err
	}

	api := slack.New(apiKey)

	messageBlocks, err := buildMessageBlocks(logDetails, logMessage, githubURL)
	if err != nil {
		logger.Error("error building message blocks", logger.Err(err))
	}
	addErrorAttachment(messageBlocks, logMessage)

	_, _, err = api.PostMessageContext(ctx, channelID, slack.MsgOptionBlocks(messageBlocks...))
	if err != nil {
		logger.Error("error posting message to Slack channel",
			logger.String("channelID", channelID),
			logger.Err(err),
		)
		return err
	}

	logger.Info("message sent to Slack channel")
	return nil
}

func buildMessageBlocks(logDetails *cloudwatch.LogDetails, logMessage *cloudwatch.LogMessage, githubURL string) ([]slack.Block, error) {
	blocks := []slack.Block{
		createTextSection(fmt.Sprintf(
			"*%s* | *%s* | Account: *%s*",
			logDetails.ErrorType,
			logDetails.Region,
			logDetails.Owner,
		)),
		createTextSection("*Message*:\n```\n" + logMessage.Msg + "\n" + logMessage.Error + "\n```"),
		createContextBlock(logDetails, logMessage, githubURL),
	}

	return blocks, nil
}

func createContextBlock(logDetails *cloudwatch.LogDetails, logMessage *cloudwatch.LogMessage, gitHubURL string) slack.Block {
	callerURL := fmt.Sprintf("<%s|%s>", gitHubURL, logMessage.Caller)

	contextText := fmt.Sprintf("*Function*: %s\n*Version*: %s\n*Caller*: %s", logDetails.FunctionName, logMessage.Version, callerURL)
	if logMessage.Tools != "" {
		contextText += fmt.Sprintf("\n*Tools*: %s", logMessage.Tools)
	}
	if logMessage.Scanner != "" {
		contextText += fmt.Sprintf("\n*Scanner*: %s", logMessage.Scanner)
	}
	contextText += fmt.Sprintf("\n*Log Stream*: <%s|%s>", logDetails.LogStreamURL, logDetails.LogStreamName)

	return slack.NewContextBlock("", slack.NewTextBlockObject("mrkdwn", contextText, false, false))
}

func createTextSection(text string) slack.Block {
	return slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", text, false, false), nil, nil)
}

func addErrorAttachment(blocks []slack.Block, logMessage *cloudwatch.LogMessage) {
	if logMessage.StdError != "" || logMessage.ErrorVerbose != "" {
		attachment := slack.Attachment{Text: fmt.Sprintf("%s\n%s", logMessage.StdError, logMessage.ErrorVerbose)}
		message := slack.NewBlockMessage(blocks...)
		message.Attachments = append(message.Attachments, attachment)
	}
}
