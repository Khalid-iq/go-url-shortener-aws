variable "region" {
  description = "AWS region for the application's resources"
  type        = string
  default     = "eu-central-1"
}

variable "project" {
  description = "Project name, used in resource names and tags"
  type        = string
  default     = "url-shortener"
}

variable "github_repo" {
  description = "GitHub repository allowed to assume the CI role, in GitHub's OIDC subject format (owner@owner_id/repo@repo_id)"
  type        = string
  default     = "Khalid-iq@183847738/go-url-shortener-aws@1432310888"
}

variable "image_tag" {
  description = "Image tag for the function's first deploy; CI updates it after that"
  type        = string
  default     = "bcd73c52547155181ee7d0600c03714061fcece3"
}
