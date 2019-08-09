data "aws_iam_policy_document" "cruddyAPI-assume-role" {
  statement {
    actions = [
      "sts:AssumeRole",
    ]

    principals {
      identifiers = ["lambda.amazonaws.com"]
      type        = "Service"
    }
  }

  version = "2012-10-17"
}

data "aws_iam_policy_document" "cruddyAPI_lambda_policy" {
  statement {
    sid       = "profileDynamoDB"
    resources = ["${aws_dynamodb_table.profiles.arn}"]
    effect    = "Allow"
    actions   = ["dynamodb:PutItem", "dynamodb:GetItem"]
  }
}

resource "aws_lambda_permission" "apigateway_lambda_invoke" {
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.cruddyAPI.arn}"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_deployment.deployment_v1.execution_arn}/*"
}

resource "aws_iam_role" "cruddyAPI_lambda" {
  assume_role_policy = "${data.aws_iam_policy_document.cruddyAPI-assume-role.json}"
}

resource "aws_iam_role_policy" "intake_lambda" {
  role   = "${aws_iam_role.cruddyAPI_lambda.id}"
  policy = "${data.aws_iam_policy_document.cruddyAPI_lambda_policy.json}"
}
