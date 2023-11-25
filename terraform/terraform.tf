terraform {
  required_version = ">= 0.15"
  backend "s3" {
    bucket  = "terraform-backend"
    key     = "cloudwatch-slack-alerts"
    region  = "us-east-2"
    encrypt = true
    profile = "shared"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0.0"
    }
  }
}

provider "aws" {
  region  = local.workspace["region"]
  profile = var.pipeline ? "" : local.workspace["profile"]
  default_tags {
    tags = {
      Repository  = local.env["global"]["tags"]["Repository"]
      Workspace   = local.env["global"]["tags"]["Workspace"]
      Service     = local.env["global"]["tags"]["Service"]
      Environment = local.env["global"]["tags"]["Environment"]
    }
  }
}

variable "pipeline" {
  description = "Set true in GHA. Allows local deployment profiles to work."
  type        = bool
  default     = false
}
