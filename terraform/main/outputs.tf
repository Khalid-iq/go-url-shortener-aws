output "ecr_repository_url" {
  description = "Push images here, e.g. <url>:<git-sha>"
  value       = aws_ecr_repository.app.repository_url
}

output "github_actions_role_arn" {
  description = "Role ARN for aws-actions/configure-aws-credentials"
  value       = aws_iam_role.github_actions.arn
}

output "github_deploy_role_arn" {
  description = "Role ARN for the deploy job (main branch only)"
  value       = aws_iam_role.github_deploy.arn
}

output "function_url" {
  description = "Public HTTPS URL of the app"
  value       = aws_lambda_function_url.app.function_url
}

output "name_servers" {
  description = "Set these as Custom DNS nameservers"
  value       = aws_route53_zone.main.name_servers
}

output "site_url" {
  description = "Public URL of the app"
  value       = "https://${var.domain_name}/"
}
