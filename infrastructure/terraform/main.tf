resource "aws_lambda_function" "cruddyAPI" {
    function_name = "cruddyAPI"
    handler = "cruddyAPI"
    filename = "../../lambda.zip"
    source_code_hash = "${base64sha256(file("../../lambda.zip"))}"
    runtime          = "go1.x"
    timeout          = "30"
    role = "${aws_iam_role.lambda_exec_role.arn}"
}