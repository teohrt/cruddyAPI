resource "aws_lambda_function" "cruddyAPI" {
  function_name    = "cruddyAPI"
  handler          = "cruddyAPI"
  filename         = "../../lambda.zip"
  source_code_hash = filebase64sha256("../../lambda.zip")
  runtime          = "go1.x"
  timeout          = "30"
  role             = "${aws_iam_role.cruddyAPI_lambda.arn}"

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = {
      LOG_LEVEL            = "debug"
      DYNAMODB_TABLE_NAME  = "${aws_dynamodb_table.profiles.name}"
      AWS_SESSION_REGION   = "${var.region}"
      AWS_SESSION_ENDPOINT = "${var.session_endpoint}"
    }
  }
}

resource "aws_dynamodb_table" "profiles" {
  name           = "profiles"
  stream_enabled = "false"

  hash_key       = "id"
  read_capacity  = 5
  write_capacity = 5

  attribute {
    name = "id"
    type = "S"
  }

  billing_mode = "PROVISIONED"

  server_side_encryption {
    enabled = true
  }
}
