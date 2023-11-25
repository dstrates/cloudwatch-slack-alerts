locals {
  env = {
    global = {
      tags = {
        Repository  = "cloudwatch-slack-alerts"
        Workspace   = terraform.workspace
        Service     = "Platform"
        Environment = terraform.workspace
      }
    }

    defaults = {
      release_version                = "0.0.1"
      enable_cloudwatch_slack_alerts = true
      role_name                      = "iam-role-cloudwatch-slack-alerts"
      lambda_function_names = [
        "cloudwatch-slack-alerts",
        # Add more lambda function names as needed to monitor for errors
      ]
      cloudwatch_lambda_errors_filter_pattern = <<EOT
    {
      ($.level = "error" || $.level = "ERROR") ||
      ($.message = "Task timed out") ||
      ($.message = "Error: Runtime exited") ||
      ($.message = "panic:") ||
      ($.message = "NetworkError") ||
      ($.message = "OutOfMemoryError") ||
      ($.message = "AccessDeniedException") ||
      ($.message = "ResourceNotFoundException") ||
      ($.message = "InvalidParameterException") ||
      ($.message = "InvalidRequestException") ||
    }
    EOT
    }

    prod-use2 = {
      region                     = "us-east-2"
      profile                    = "prod"
      account_id                 = "PLACEHOLDER"
      s3_key_prefix              = "PLACEHOLDER"
      s3_existing_package_bucket = "PLACEHOLDER"
    }
  }

  workspace = merge(local.env[defaults], local.env[terraform.workspace])
}
