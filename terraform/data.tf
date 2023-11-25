data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

data "aws_iam_role" "alerter" {
  name = "iam-role-cloudwatch-slack-alerts"
}
