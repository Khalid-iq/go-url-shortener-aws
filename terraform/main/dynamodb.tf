resource "aws_dynamodb_table" "links" {
  name         = "${var.project}-links"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "code"

  attribute {
    name = "code"
    type = "S"
  }

  # Continuous backups: restore the table to any second in the last 35 days.
  point_in_time_recovery {
    enabled = true
  }

  # Blocks deletion from the console, CLI, or Terraform until turned off.
  deletion_protection_enabled = true
}

data "aws_iam_policy_document" "lambda_dynamodb" {
  statement {
    actions   = ["dynamodb:GetItem", "dynamodb:PutItem"]
    resources = [aws_dynamodb_table.links.arn]
  }
}

resource "aws_iam_role_policy" "lambda_dynamodb" {
  name   = "links-table"
  role   = aws_iam_role.lambda_exec.id
  policy = data.aws_iam_policy_document.lambda_dynamodb.json
}
