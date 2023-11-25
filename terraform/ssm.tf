resource "aws_ssm_parameter" "slack_key" {
  lifecycle {
    create_before_destroy = true
    ignore_changes        = [value]
  }
  name        = "/slack/cloudwatch-alerts/key"
  description = "Slack key for CloudWatch error alerts"
  type        = "SecureString"
  value       = "PLACEHOLDER"
}

resource "aws_ssm_parameter" "slack_channel_mapping" {
  lifecycle {
    create_before_destroy = true
  }
  name        = "/slack/cloudwatch-alerts/channel-map"
  description = "Map of service name to channel ID for Slack alerts"
  type        = "String"
  value       = <<EOF
{
  "service-1": "PLACEHOLDER",
  "service-2": "PLACEHOLDER",
}
EOF
}

resource "aws_ssm_parameter" "default_slack_channel_id" {
  lifecycle {
    create_before_destroy = true
  }
  name        = "/slack/cloudwatch-alerts/default-channel-id"
  description = "Default channel ID for Slack alerts"
  type        = "String"
  value       = "PLACEHOLDER"
}
