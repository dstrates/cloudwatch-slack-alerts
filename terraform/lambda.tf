module "lambda_alerter" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.2.1"

  function_name                           = "cloudwatch-slack-alerts"
  description                             = "Sends alerts to Slack on CloudWatch error events"
  handler                                 = "main"
  runtime                                 = "go1.x"
  ephemeral_storage_size                  = 512 # Min 512 MB to 10,240 MB (10 GB)
  memory_size                             = 128 # Min 128 MB to 10,240 MB (10 GB)
  architectures                           = ["x86_64"]
  create_role                             = false
  lambda_role                             = try(module.iam_role_alerter[0].iam_role_arn, data.aws_iam_role.alerter.arn)
  create_package                          = false
  ignore_source_code_hash                 = false
  publish                                 = true
  create_current_version_allowed_triggers = true
  timeout                                 = 3
  s3_existing_package = {
    bucket = local.workspace.s3_existing_package_bucket
    key    = "${local.workspace.s3_key_prefix}/${local.workspace.release_version}/alerter.zip"
  }
  environment_variables = {
    LOG_LEVEL                            = "error"
    VERSION                              = local.workspace.release_version
    SLACK_KEY_PARAMETER_NAME             = aws_ssm_parameter.slack_key.name
    DEFAULT_SLACK_CHANNEL_PARAMETER_NAME = aws_ssm_parameter.default_slack_channel_id.name
  }
  allowed_triggers = {
    CloudWatch = {
      principal  = "logs.${local.workspace.region}.amazonaws.com"
      source_arn = "arn:aws:logs:${local.workspace.region}:${local.workspace.account_id}:log-group:/aws/lambda/*:*"
    },
  }
}
