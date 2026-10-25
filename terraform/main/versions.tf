terraform {
  required_version = ">= 1.10"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }


  backend "s3" {
    bucket       = "url-shortener-tfstate-541271671872-eu-central-1"
    key          = "main/terraform.tfstate"
    region       = "eu-central-1"
    encrypt      = true
    use_lockfile = true # S3-native locking, no DynamoDB table needed
  }
}

provider "aws" {
  region = var.region

  default_tags {
    tags = {
      Project   = var.project
      ManagedBy = "terraform"
      Stack     = "main"
    }
  }
}
