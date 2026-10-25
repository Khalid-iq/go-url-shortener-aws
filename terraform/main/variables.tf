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
  default     = "f983ee075ba1b052fa947f67123cf5c76fcee226"
}

variable "domain_name" {
  description = "Public domain for the shortener"
  type        = string
  default     = "kqlinks.com"
}
