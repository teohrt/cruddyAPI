resource "aws_api_gateway_rest_api" "cruddyAPI" {
  name = "cruddyAPI-${terraform.workspace}-gateway"
  description = "This is my API for demonstration purposes"
}

resource "aws_api_gateway_deployment" "deployment_v1" {
  rest_api_id = "${aws_api_gateway_rest_api.cruddyAPI.id}"
  stage_name  = "v1"
}

resource "aws_api_gateway_resource" "system" {
  rest_api_id = "${aws_api_gateway_rest_api.cruddyAPI.id}"
  parent_id   = "${aws_api_gateway_rest_api.cruddyAPI.root_resource_id}"
  path_part   = "system"
}

resource "aws_api_gateway_resource" "ping" {
  rest_api_id = "${aws_api_gateway_rest_api.cruddyAPI.id}"
  parent_id   = "${aws_api_gateway_resource.system.id}"
  path_part   = "ping"
}

resource "aws_api_gateway_resource" "cruddyAPI" {
  parent_id   = "${aws_api_gateway_rest_api.cruddyAPI.root_resource_id}"
  path_part   = "cruddyAPI"
  rest_api_id = "${aws_api_gateway_rest_api.cruddyAPI.id}"
}

resource "aws_api_gateway_method" "cruddyAPI" {
  http_method      = "POST"
  resource_id      = "${aws_api_gateway_resource.cruddyAPI.id}"
  rest_api_id      = "${aws_api_gateway_rest_api.cruddyAPI.id}"
  authorization    = "NONE"
  api_key_required = true
}

resource "aws_api_gateway_integration" "cruddyAPI" {
  http_method             = "POST"
  resource_id             = "${aws_api_gateway_resource.cruddyAPI.id}"
  rest_api_id             = "${aws_api_gateway_rest_api.cruddyAPI.id}"
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.cruddyAPI.invoke_arn}"
}

resource "aws_api_gateway_method" "ping" {
  rest_api_id      = "${aws_api_gateway_rest_api.cruddyAPI.id}"
  resource_id      = "${aws_api_gateway_resource.ping.id}"
  http_method      = "GET"
  authorization    = "NONE"
  api_key_required = true
}

resource "aws_api_gateway_integration" "ping" {
  rest_api_id          = "${aws_api_gateway_rest_api.cruddyAPI.id}"
  resource_id          = "${aws_api_gateway_method.ping.resource_id}"
  http_method          = "${aws_api_gateway_method.ping.http_method}"
  type                 = "MOCK"
  timeout_milliseconds = 1000

  request_templates {
    "application/json" = "{\"statusCode\": 200}"
  }
}

resource "aws_api_gateway_integration_response" "ping" {
  rest_api_id = "${aws_api_gateway_rest_api.cruddyAPI.id}"
  resource_id = "${aws_api_gateway_resource.ping.id}"
  http_method = "${aws_api_gateway_method.ping.http_method}"
  status_code = "${aws_api_gateway_method_response.ping.status_code}"

  depends_on = [
    "aws_api_gateway_integration.ping",
    "aws_api_gateway_method_response.ping",
  ]
}

resource "aws_api_gateway_method_response" "ping" {
  rest_api_id = "${aws_api_gateway_rest_api.cruddyAPI.id}"
  resource_id = "${aws_api_gateway_resource.ping.id}"
  http_method = "${aws_api_gateway_method.ping.http_method}"
  status_code = 200

  response_models = {
    "application/json" = "Empty"
  }

  depends_on = [
    "aws_api_gateway_method.ping",
  ]
}