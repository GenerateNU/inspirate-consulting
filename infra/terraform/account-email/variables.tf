variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "name_prefix" {
  type    = string
  default = "inspirate"
}

variable "sender_email" {
  type = string
}

variable "app_url" {
  type = string
}

variable "claim_token_secret" {
  type      = string
  sensitive = true
}

variable "lambda_zip_path" {
  type    = string
  default = "${path.module}/build/account-email.zip"
}

variable "lambda_timeout_seconds" {
  type    = number
  default = 15
}

variable "visibility_timeout_seconds" {
  type    = number
  default = 90
}

variable "message_retention_seconds" {
  type    = number
  default = 86400
}
