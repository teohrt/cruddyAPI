GOBUILD=go build
TERRAFORM=cd infrastructure/terraform; terraform
LAMBDA_PARAMS=env GOOS=linux GOARCH=amd64

BINARY_NAME=cruddyAPI
LOCAL_SERVER_PORT=8002
LOCAL_DB_PORT=8000

BASE_ENV_VALS := SERVER_PORT="$(LOCAL_SERVER_PORT)" \
		LOG_LEVEL="debug" \
		DYNAMODB_TABLE_NAME="profiles" \
		AWS_SESSION_REGION="us-east-2" \
		AWS_SESSION_ENDPOINT="http://localhost:$(LOCAL_DB_PORT)"

# General commands
clean:
	rm -f $(BINARY_NAME) lambda.zip

test-ci:
	go test ./... -race -cover -v 2>&1

# Local development commands
run-locally: build
	$(BASE_ENV_VALS) ./$(BINARY_NAME)

build:
	$(GOBUILD)

db-start:
	docker run -p $(LOCAL_DB_PORT):$(LOCAL_DB_PORT) amazon/dynamodb-local

db-table-init:
	cd dbclient; aws dynamodb create-table --cli-input-json file://profiles_table.json --endpoint-url http://localhost:$(LOCAL_DB_PORT)

# Infrastructure commands
deploy: clean build-lambda zip-lambda init
	$(TERRAFORM) apply

destroy:
	$(TERRAFORM) destroy

init:
	$(TERRAFORM) init || exit 1

build-lambda:
	make test-ci || exit 1
	$(LAMBDA_PARAMS) $(GOBUILD)

zip-lambda:
	zip -j lambda.zip $(BINARY_NAME)