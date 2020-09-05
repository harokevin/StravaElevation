provider "aws" {
  profile = "default"
  region  = "us-west-2"
}

resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

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

resource "aws_lambda_function" "strava_elevation_TF_demo" {
  filename      = "StravaElevation.zip"
  function_name = "strava_elevation_TF_demo"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "StravaElevation"

  source_code_hash = filebase64sha256("StravaElevation.zip")

  runtime = "go1.x"

  environment {
    variables = {
      client_id = var.client_id
      client_secret = var.client_secret
      refresh_token = var.refresh_token
    }
  }
}

resource "aws_cloudwatch_event_rule" "during_day" {
  name        = "during_day"
  description = "Run the StravaElevation lambda during the day"
  schedule_expression = "cron(0 20,21,22,23,0,1,2,3,4 * * ? *)"
}

resource "aws_cloudwatch_event_target" "strava_elevation" {
  rule      = aws_cloudwatch_event_rule.during_day.name
  target_id = "lambda"
  arn       = aws_lambda_function.strava_elevation_TF_demo.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_run_strava_elevation" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.strava_elevation_TF_demo.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.during_day.arn
}