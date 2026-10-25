resource "aws_route53_zone" "main" {
  name = var.domain_name


  lifecycle {
    prevent_destroy = true
  }
}
