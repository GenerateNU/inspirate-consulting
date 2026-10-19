variable "aws_region" {
  description = "AWS region the queue and Lambda are created in. The existing task-notification queue lives in us-east-1, so this stays consistent with it."
  type        = string
  default     = "us-east-1"
}

variable "name_prefix" {
  description = "Prefix applied to every resource name, matching the existing inspirate-email-notification-queue convention."
  type        = string
  default     = "inspirate"
}

variable "sender_email" {
  description = "Verified SES identity the Lambda sends from (SENDER_EMAIL env var). Must already be a verified SES identity; SES is in sandbox mode, so recipients must be verified too until production access is granted."
  type        = string
}

variable "app_url" {
  description = "Base URL of the frontend (no trailing slash), used to build the claim-password link embedded in the email, e.g. https://inspirate.org"
  type        = string
}

variable "claim_token_secret" {
  description = "HMAC secret used to sign/verify claim-password tokens. Must match the CLAIM_TOKEN_SECRET configured on the backend API (see backend/internal/claimtoken) so tokens minted by the Lambda verify there."
  type        = string
  sensitive   = true
}

variable "lambda_zip_path" {
  description = "Path to the zipped Lambda deployment package, built from backend/cmd/lambda/account-email. See README.md in this directory for the build command."
  type        = string
  default     = "${path.module}/build/account-email.zip"
}

variable "lambda_timeout_seconds" {
  description = "Lambda timeout. 15s to allow for an SES call after formatting; the existing email-notification lambda uses 10s with a simpler (text-only) body."
  type        = number
  default     = 15
}

variable "visibility_timeout_seconds" {
  description = "SQS visibility timeout. Kept comfortably above the Lambda timeout so a slow invocation doesn't let SQS redeliver the same message to a second concurrent invocation."
  type        = number
  default     = 90
}

variable "message_retention_seconds" {
  description = "How long an unconsumed message stays in the queue before SQS drops it."
  type        = number
  default     = 86400 # 1 day
}
