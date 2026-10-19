terraform {
  required_version = ">= 1.5"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.4"
    }
  }

  # No backend configured here on purpose — this module doesn't assume a
  # remote state bucket/table exists yet. Add a `backend "s3" {}` block (and
  # run `terraform init -migrate-state`) once one is provisioned.
}

provider "aws" {
  region = var.aws_region
}

# ---------------------------------------------------------------------------
# Queue
# ---------------------------------------------------------------------------

# Messages that fail processing 5 times land here instead of retrying forever
# or silently vanishing after message_retention_seconds.
resource "aws_sqs_queue" "account_email_dlq" {
  name                      = "${var.name_prefix}-account-email-queue-dlq"
  message_retention_seconds = 1209600 # 14 days, the SQS maximum
}

resource "aws_sqs_queue" "account_email" {
  name                       = "${var.name_prefix}-account-email-queue"
  visibility_timeout_seconds = var.visibility_timeout_seconds
  message_retention_seconds  = var.message_retention_seconds

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.account_email_dlq.arn
    maxReceiveCount     = 5
  })
}

# ---------------------------------------------------------------------------
# Lambda
# ---------------------------------------------------------------------------

data "archive_file" "account_email_lambda" {
  type        = "zip"
  source_file = var.lambda_zip_path
  output_path = "${path.module}/build/account-email.zip"
}

resource "aws_iam_role" "account_email_lambda" {
  name = "${var.name_prefix}-account-email-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

# Covers CloudWatch Logs (CreateLogGroup/Stream, PutLogEvents) so invocations
# are debuggable, same as the existing email-notification lambda relies on.
resource "aws_iam_role_policy_attachment" "basic_execution" {
  role       = aws_iam_role.account_email_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy" "account_email_lambda" {
  name = "${var.name_prefix}-account-email-lambda-policy"
  role = aws_iam_role.account_email_lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "ConsumeAccountEmailQueue"
        Effect = "Allow"
        Action = [
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes",
        ]
        Resource = aws_sqs_queue.account_email.arn
      },
      {
        # Scoped to "*" because SES doesn't let you restrict SendEmail to a
        # single verified identity via resource ARN in all regions/setups;
        # tighten to the SES identity ARN once the sender domain/email is
        # consistently provisioned via Terraform too.
        Sid      = "SendAccountEmail"
        Effect   = "Allow"
        Action   = ["ses:SendEmail", "ses:SendRawEmail"]
        Resource = "*"
      },
    ]
  })
}

resource "aws_lambda_function" "account_email" {
  function_name = "${var.name_prefix}-account-email-lambda"
  role          = aws_iam_role.account_email_lambda.arn

  filename         = data.archive_file.account_email_lambda.output_path
  source_code_hash = data.archive_file.account_email_lambda.output_base64sha256

  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  timeout       = var.lambda_timeout_seconds

  environment {
    variables = {
      SENDER_EMAIL       = var.sender_email
      APP_URL            = var.app_url
      CLAIM_TOKEN_SECRET = var.claim_token_secret
    }
  }
}

resource "aws_lambda_event_source_mapping" "account_email" {
  event_source_arn = aws_sqs_queue.account_email.arn
  function_name    = aws_lambda_function.account_email.arn
  batch_size       = 10
}
