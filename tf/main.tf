terraform {
  required_version = "~> 0.15.0"

  backend "s3" {
    bucket = "bscfeestfstate"
    region = "us-east-1"
    key    = "network/terraform.tfstate"
  }
}

provider "aws" {
  region = var.aws_region
  allowed_account_ids = [
    var.aws_account_id
  ]
}

locals {
  project      = "bsf-fees"
  s3_origin_id = "bscfeesOrigin"
  tags = {
    Project   = "BSC Fees"
    Terraform = "true"
    Owner     = "katsadim"
  }
}

data "terraform_remote_state" "s3" {
  backend = "s3"
  config = {
    bucket  = "bscfeestfstate"
    key     = "network/terraform.tfstate"
    region  = var.aws_region
    encrypt = true
  }
}

############### Remote State ################

resource "aws_s3_bucket" "bsc_fees_tf_state" {
  bucket        = "bscfeestfstate"
  acl           = "private"
  force_destroy = true

  tags = local.tags
}

############### Lambda ################

resource "aws_s3_bucket" "bsc_fees_api" {
  bucket        = var.api_s3_bucket
  acl           = "private"
  force_destroy = true

  tags = local.tags
}

// This is a dummy executable that needs to be uploaded for the lambda to work.
// It is meant to be overwritten
resource "aws_s3_bucket_object" "file_upload" {
  bucket = aws_s3_bucket.bsc_fees_api.id
  key    = var.api_zip_name
  source = "bin/${var.api_zip_name}"
}

resource "aws_lambda_function" "bsc_fees" {
  function_name = var.lambda_function_name
  handler       = "bsc-fees"
  role          = aws_iam_role.iam_for_lambda.arn
  runtime       = "go1.x"
  memory_size   = 128
  timeout       = 30
  s3_bucket     = aws_s3_bucket.bsc_fees_api.id
  s3_key        = var.api_zip_name

  environment {
    variables = {
      GENERAL_ENV    = "prod"
      BINANCE_APIKEY = var.binance_api_key
      BSC_APIKEY     = var.bscscan_api_key
      ETH_APIKEY     = var.ethscan_api_key
    }
  }

  tags = local.tags
  depends_on = [
    aws_iam_role_policy_attachment.lambda_logs,
    aws_cloudwatch_log_group.lambda_log_group,
  ]
}

resource "aws_iam_role" "iam_for_lambda" {
  name               = "serverless_example_lambda"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

}

resource "aws_cloudwatch_log_group" "lambda_log_group" {
  name              = "/aws/lambda/${var.lambda_function_name}"
  retention_in_days = var.log_retention_in_days

  tags = local.tags
}

# See also the following AWS managed policy: AWSLambdaBasicExecutionRole
resource "aws_iam_policy" "lambda_logging" {
  name        = "lambda_logging"
  path        = "/"
  description = "IAM policy for logging from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.lambda_logging.arn
}

resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.bsc_fees.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.aws_region}:${var.aws_account_id}:${aws_api_gateway_rest_api.rest_api.id}/*/${aws_api_gateway_method.bsc_fees.http_method}${aws_api_gateway_resource.bsc_fees.path}"
}

################### API Gateway ##########################

resource "aws_api_gateway_rest_api" "rest_api" {
  name        = "Public Api"
  description = "BSC Fees ftw"
  tags        = local.tags
}

resource "aws_api_gateway_resource" "bsc_fees" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  parent_id   = aws_api_gateway_rest_api.rest_api.root_resource_id
  path_part   = "fees"
}

resource "aws_api_gateway_method" "bsc_fees" {
  rest_api_id          = aws_api_gateway_rest_api.rest_api.id
  resource_id          = aws_api_gateway_resource.bsc_fees.id
  http_method          = "GET"
  authorization        = "NONE"
  request_validator_id = aws_api_gateway_request_validator.bsc_fees.id
  request_parameters = {
    "method.request.querystring.address" = true
  }
}

resource "aws_api_gateway_method" "cors" {
  rest_api_id   = aws_api_gateway_rest_api.rest_api.id
  resource_id   = aws_api_gateway_resource.bsc_fees.id
  http_method   = "OPTIONS"
  authorization = "NONE"
}

resource "aws_api_gateway_request_validator" "bsc_fees" {
  name                        = "Validate Query Params"
  rest_api_id                 = aws_api_gateway_rest_api.rest_api.id
  validate_request_body       = false
  validate_request_parameters = true
}

resource "aws_api_gateway_integration" "lambda" {
  rest_api_id             = aws_api_gateway_rest_api.rest_api.id
  resource_id             = aws_api_gateway_method.bsc_fees.resource_id
  http_method             = aws_api_gateway_method.bsc_fees.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.bsc_fees.invoke_arn

  request_parameters = {
    "integration.request.querystring.address" = "method.request.querystring.address"
  }
}

resource "aws_api_gateway_integration" "cors" {
  rest_api_id             = aws_api_gateway_rest_api.rest_api.id
  resource_id             = aws_api_gateway_method.cors.resource_id
  http_method             = aws_api_gateway_method.cors.http_method
  integration_http_method = "OPTIONS"
  type                    = "MOCK"
}

resource "aws_api_gateway_method_response" "lambda_200" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  resource_id = aws_api_gateway_method.bsc_fees.resource_id
  http_method = aws_api_gateway_method.bsc_fees.http_method
  status_code = "200"

  depends_on = [
    aws_api_gateway_integration.lambda,
  ]
}

resource "aws_api_gateway_method_response" "lambda_400" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  resource_id = aws_api_gateway_method.bsc_fees.resource_id
  http_method = aws_api_gateway_method.bsc_fees.http_method
  status_code = "400"

  depends_on = [
    aws_api_gateway_integration.lambda,
  ]
}

resource "aws_api_gateway_method_response" "lambda_500" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  resource_id = aws_api_gateway_method.bsc_fees.resource_id
  http_method = aws_api_gateway_method.bsc_fees.http_method
  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin"  = true,
    "method.response.header.Access-Control-Allow-Methods" = true,
    "method.response.header.Access-Control-Allow-Headers" = true,
  }
  status_code = "500"

  depends_on = [
    aws_api_gateway_integration.lambda,
  ]
}

resource "aws_api_gateway_integration_response" "cors_integration" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  resource_id = aws_api_gateway_method.cors.resource_id
  http_method = aws_api_gateway_method.cors.http_method
  status_code = "200"
  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin"  = "'https://bscfees.com'",
    "method.response.header.Access-Control-Allow-Methods" = "'GET OPTIONS'",
    "method.response.header.Access-Control-Allow-Headers" = "'Content-Type'",
  }

  depends_on = [
    aws_api_gateway_integration.cors,
  ]
}

resource "aws_api_gateway_method_response" "cors" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  resource_id = aws_api_gateway_method.cors.resource_id
  http_method = aws_api_gateway_method.cors.http_method
  status_code = "200"
  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin"  = true,
    "method.response.header.Access-Control-Allow-Methods" = true,
    "method.response.header.Access-Control-Allow-Headers" = true,
  }

  depends_on = [
    aws_api_gateway_integration.cors,
  ]
}

############## API Gateway Settings #####################

resource "aws_api_gateway_method_settings" "settings" {
  method_path = "*/*"
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  stage_name  = aws_api_gateway_stage.stage.stage_name

  settings {
    metrics_enabled        = true
    logging_level          = "INFO"
    data_trace_enabled     = true
    throttling_rate_limit  = 5
    throttling_burst_limit = 3
  }
}

data "aws_iam_policy_document" "cloudwatch_apigateway_assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      identifiers = ["apigateway.amazonaws.com"]
      type        = "Service"
    }
  }
}

resource "aws_iam_role" "cloudwatch_apigateway_role" {
  name               = "${local.project}-cloudwatch-apigateway-role"
  assume_role_policy = data.aws_iam_policy_document.cloudwatch_apigateway_assume_role_policy.json
}

resource "aws_iam_role_policy_attachment" "cloudwatch_apigateway_policy" {
  role       = aws_iam_role.cloudwatch_apigateway_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs"
}

resource "aws_api_gateway_account" "account" {
  cloudwatch_role_arn = aws_iam_role.cloudwatch_apigateway_role.arn
}

resource "aws_cloudwatch_log_group" "rest_api_execution_log" {
  name              = "Api-GW-Execution-Logs-${aws_api_gateway_rest_api.rest_api.id}"
  retention_in_days = var.log_retention_in_days

  tags = local.tags
}

############## API Gateway Deployment #####################

resource "aws_api_gateway_stage" "stage" {
  deployment_id = aws_api_gateway_deployment.deployment.id
  rest_api_id   = aws_api_gateway_rest_api.rest_api.id
  stage_name    = var.api_gateway_stage_name
  description   = "Public API ${timestamp()}"
  tags          = local.tags

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.rest_api_access_log.arn
    format          = file("text/rest_api_gateway_access_log_format.json")
  }

  depends_on = [
  aws_cloudwatch_log_group.rest_api_access_log]
}

resource "aws_api_gateway_deployment" "deployment" {
  rest_api_id       = aws_api_gateway_rest_api.rest_api.id
  description       = "Bsc-Fees API"
  stage_description = "Api Gateway that support BscFees"
  # use the following if you want to invalidate the dpeloyment
  # stage_description = "Deployed at ${timestamp()}"

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [
    aws_api_gateway_rest_api.rest_api,
    aws_api_gateway_integration.lambda,
    aws_api_gateway_integration.cors,
    aws_api_gateway_method.bsc_fees,
    aws_api_gateway_method.cors,
  ]
}

resource "aws_cloudwatch_log_group" "rest_api_access_log" {
  name              = "Api-GW-Access-Logs-${aws_api_gateway_rest_api.rest_api.id}"
  retention_in_days = var.log_retention_in_days

  tags = local.tags
}

#####################  API GW Custom Domain Name  #####################

resource "aws_api_gateway_domain_name" "api_gateway_domain_name" {
  certificate_arn = aws_acm_certificate_validation.bscfees.certificate_arn
  domain_name     = "api.bscfees.com"
  security_policy = "TLS_1_2"

  tags = local.tags
}

resource "aws_route53_record" "apigw" {
  name    = aws_api_gateway_domain_name.api_gateway_domain_name.domain_name
  type    = "A"
  zone_id = aws_route53_zone.bscfees.id

  alias {
    evaluate_target_health = true
    name                   = aws_api_gateway_domain_name.api_gateway_domain_name.cloudfront_domain_name
    zone_id                = aws_api_gateway_domain_name.api_gateway_domain_name.cloudfront_zone_id
  }
}

resource "aws_api_gateway_base_path_mapping" "bscfees_mapping" {
  api_id      = aws_api_gateway_rest_api.rest_api.id
  stage_name  = aws_api_gateway_stage.stage.stage_name
  domain_name = aws_api_gateway_domain_name.api_gateway_domain_name.domain_name
}

#####################  Certificate  #####################

resource "aws_acm_certificate" "bscfees" {
  domain_name               = "bscfees.com"
  subject_alternative_names = ["www.bscfees.com", "api.bscfees.com"]
  validation_method         = "DNS"

  tags = local.tags

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route53_zone" "bscfees" {
  name = "bscfees.com"

  tags = local.tags
}

resource "aws_route53_record" "bscfees" {
  for_each = {
    for dvo in aws_acm_certificate.bscfees.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = aws_route53_zone.bscfees.zone_id
}

resource "aws_acm_certificate_validation" "bscfees" {
  certificate_arn         = aws_acm_certificate.bscfees.arn
  validation_record_fqdns = [for record in aws_route53_record.bscfees : record.fqdn]
}

#####################    Web Cloudfront Distribution   #####################

resource "aws_s3_bucket" "bsc_fees_web" {
  bucket        = var.web_s3_bucket
  acl           = "public-read"
  force_destroy = true

  website {
    index_document = "index.html"
    error_document = "404.html"
  }

  tags = local.tags
}

resource "aws_cloudfront_distribution" "s3_distribution" {
  origin {
    domain_name = aws_s3_bucket.bsc_fees_web.bucket_regional_domain_name
    origin_id   = local.s3_origin_id

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.origin_access_identity.cloudfront_access_identity_path
    }
  }

  aliases = ["www.bscfees.com", "bscfees.com"]

  enabled             = true
  is_ipv6_enabled     = true
  comment             = "BscFees static site distribution"
  default_root_object = "index.html"

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = local.s3_origin_id

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  price_class = "PriceClass_200"

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn = aws_acm_certificate.bscfees.arn
    ssl_support_method  = "sni-only"
  }

  tags = local.tags
}

resource "aws_cloudfront_origin_access_identity" "origin_access_identity" {
  comment = "Some comment"
}

resource "aws_route53_record" "web_ipv4" {
  name    = "www.bscfees.com"
  type    = "A"
  zone_id = aws_route53_zone.bscfees.id

  alias {
    evaluate_target_health = true
    name                   = aws_cloudfront_distribution.s3_distribution.domain_name
    zone_id                = aws_cloudfront_distribution.s3_distribution.hosted_zone_id
  }
}

resource "aws_route53_record" "web_ipv6" {
  name    = "www.bscfees.com"
  type    = "AAAA"
  zone_id = aws_route53_zone.bscfees.id

  alias {
    evaluate_target_health = true
    name                   = aws_cloudfront_distribution.s3_distribution.domain_name
    zone_id                = aws_cloudfront_distribution.s3_distribution.hosted_zone_id
  }
}

resource "aws_route53_record" "redirect_non_www_to_www" {
  name    = "bscfees.com"
  type    = "A"
  zone_id = aws_route53_zone.bscfees.id

  alias {
    evaluate_target_health = true
    name                   = "www.bscfees.com"
    zone_id                = aws_route53_zone.bscfees.id
  }
}
