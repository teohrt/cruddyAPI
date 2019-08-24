resource "aws_api_gateway_rest_api" "cruddyAPI" {
  name        = "cruddyAPI-${terraform.workspace}-gateway"
  description = "This is my API for demonstration purposes"
}

resource "aws_api_gateway_deployment" "deployment_v1" {
  rest_api_id = "${aws_api_gateway_rest_api.cruddyAPI.id}"
  stage_name  = "api"

  depends_on = [
    "aws_api_gateway_integration.cruddyAPI",
  ]
}

resource "aws_api_gateway_resource" "cruddyAPIProxy" {
  parent_id   = "${aws_api_gateway_rest_api.cruddyAPI.root_resource_id}"
  path_part   = "{proxy+}"
  rest_api_id = "${aws_api_gateway_rest_api.cruddyAPI.id}"
}

resource "aws_api_gateway_method" "cruddyAPI" {
  http_method      = "ANY"
  resource_id      = "${aws_api_gateway_resource.cruddyAPIProxy.id}"
  rest_api_id      = "${aws_api_gateway_rest_api.cruddyAPI.id}"
  authorization    = "NONE"
  api_key_required = false

  request_parameters = {
    "method.request.path.proxy" = true
  }
}

resource "aws_api_gateway_integration" "cruddyAPI" {
  http_method             = "${aws_api_gateway_method.cruddyAPI.http_method}"
  resource_id             = "${aws_api_gateway_resource.cruddyAPIProxy.id}"
  rest_api_id             = "${aws_api_gateway_rest_api.cruddyAPI.id}"
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.cruddyAPI.invoke_arn}"

  depends_on = [
    "aws_api_gateway_method.cruddyAPI",
  ]
}
