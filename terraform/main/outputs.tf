output "ecr_repository_url" {
  description = "Push images here, e.g. <url>:<git-sha>"
  value       = aws_ecr_repository.app.repository_url
}
