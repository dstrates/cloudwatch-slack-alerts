# Example resource that needs to be created with each Lambda function to be monitored
# This isn't required for the alerter lambda to function
# See locals for lambda_function_name declaration method

resource "aws_cloudwatch_log_subscription_filter" "lambda_errors" {
  count = local.workspace.enable_cloudwatch_slack_alerts ? length(local.workspace.lambda_function_names) : 0

  depends_on = [
    module.lambda_alerter,
  ]

  name            = "alerter-${local.workspace.lambda_function_names[count.index]}"
  log_group_name  = "/aws/lambda/${local.workspace.lambda_function_names[count.index]}"
  filter_pattern  = local.workspace.cloudwatch_lambda_errors_filter_pattern
  destination_arn = data.aws_lambda_function.lambda_alerts.arn
}
