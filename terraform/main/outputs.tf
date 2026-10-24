output "ecr_repository_url" {
  description = "Push images here, e.g. <url>:<git-sha>"
  value       = aws_ecr_repository.app.repository_url
}

output "github_actions_role_arn" {
  description = "Role ARN for aws-actions/configure-aws-credentials"
  value       = aws_iam_role.github_actions.arn
}
