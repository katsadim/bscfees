variable "api_s3_bucket" {
  description = "The S3 bucket that hosts the lambda function code"
  type        = string
}

variable "api_gateway_stage_name" {
  description = "The API Gateway stage name"
  type        = string
}

variable "api_zip_name" {
  description = "The name of the zip fiel that hosts the binary API application"
  type        = string
}

variable "aws_account_id" {
  description = "The id of the AWS Account to apply changes to."
  type        = string
}

variable "aws_region" {
  description = "The region of our accounts."
  type        = string
}

variable "bscscan_api_key" {
  description = "Binance Smart Chain Developer API key (bscscan.com)"
  type        = string
}

variable "binance_api_key" {
  description = "Binance Developer API key (binance.com)"
  type        = string
}

variable "web_s3_bucket" {
  description = "The S3 bucket that hosts the web code"
  type        = string
}

variable "lambda_function_name" {
  description = "The name of the lambda funciton (duh)."
  type        = string
}

variable "log_retention_in_days" {
  description = "The number of days we are going to keep logs for"
  type        = number
}
