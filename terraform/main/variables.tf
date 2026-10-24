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
  description = "GitHub repository allowed to assume the CI role (owner/name)"
  type        = string
  default     = "Khalid-iq/go-url-shortener-aws"
}
