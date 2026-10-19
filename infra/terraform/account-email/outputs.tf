output "queue_url" {
  value = aws_sqs_queue.account_email.id
}

output "queue_arn" {
  value = aws_sqs_queue.account_email.arn
}

output "dead_letter_queue_url" {
  value = aws_sqs_queue.account_email_dlq.id
}

output "lambda_function_name" {
  value = aws_lambda_function.account_email.function_name
}

output "lambda_function_arn" {
  value = aws_lambda_function.account_email.arn
}
