output "state_bucket" {
  description = "Use this as the bucket name in the main config's backend block"
  value       = aws_s3_bucket.tfstate.bucket
}
