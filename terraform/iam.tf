module "iam_role_alerter_lambda" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-assumable-role"
  version = "5.34.0"
  count   = try(local.workspace.create_iam_role_alerter, true) ? 1 : 0

  role_name             = "iam-role-cloudwatch-slack-alerts"
  create_role           = true
  role_requires_mfa     = false
  trusted_role_services = ["lambda.amazonaws.com"]
  custom_role_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
    module.iam_policy_alerter_lambda[0].arn
  ]
}

module "iam_policy_alerter_lambda" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-policy"
  version = "5.34.0"
  count   = try(local.workspace.create_iam_role_alerter, true) ? 1 : 0

  name   = "iam-policy-cloudwatch-slack-alerts"
  policy = data.aws_iam_policy_document.alerter.json
}

data "aws_iam_policy_document" "alerter" {
  statement {
    actions = ["ssm:GetParameter"]
    resources = [
      "${aws_ssm_parameter.slack_key.arn}",
      "${aws_ssm_parameter.slack_channel_mapping.arn}",
      "${aws_ssm_parameter.default_slack_channel_id.arn}"
    ]
  }
}
