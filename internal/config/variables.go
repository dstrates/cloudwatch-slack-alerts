package config

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/nullify-platform/logger/pkg/logger"
)

func getStringVariable(name string, defaultValue string) string {
	val := os.Getenv(name)
	if val == "" {
		return defaultValue
	}
	return val
}

func getBooleanVariable(name string, defaultValue bool) bool {
	val := os.Getenv(name)
	if val == "" {
		return defaultValue
	}
	if val == "true" {
		return true
	}
	return false
}

func getParameterStoreVariable(ctx context.Context, name string, altNames ...string) (string, error) {
	parameterName := os.Getenv(name)
	if parameterName == "" {
		// loop through alternate names until one is defined
		for _, altName := range altNames {
			parameterName = os.Getenv(altName)
			if parameterName != "" {
				break
			}
		}
	}
	if parameterName == "" {
		return "", errors.New("parameter name is empty")
	}
	cfg, err := NewAWSConfig(ctx)
	if err != nil {
		logger.Error(
			"creating new aws config",
			logger.Err(err),
		)
		return "", err
	}
	svc := ssm.NewFromConfig(cfg)
	param, err := svc.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(parameterName),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}
	return *param.Parameter.Value, nil
}
