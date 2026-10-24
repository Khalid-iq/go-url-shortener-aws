variable "region" {
  description = "AWS region for the state bucket"
  type        = string
  default     = "eu-central-1"
}

variable "project" {
  description = "Project name, used in resource names and tags"
  type        = string
  default     = "url-shortener"
}
