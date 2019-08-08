GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=cruddyAPI
LOCAL_PORT="8002"

BASE_ENV_VALS := _LAMBDA_SERVER_PORT=$(LOCAL_PORT) \
		LOG_LEVEL="debug" \
		DYNAMODB_TABLE_NAME="profiles" \
		AWS_SESSION_REGION="us-east-2" \
		AWS_SESSION_ENDPOINT="http://localhost:8000"

run: build
	$(BASE_ENV_VALS) ./$(BINARY_NAME)

build:
	$(GOBUILD)
	
clean:
	rm -f $(BINARY_NAME)

test-ci:
	go test ./... -race -cover -v 2>&1

build-lambda:
	env GOOS=linux GOARCH=amd64 $(GOBUILD)

zip:
	zip -j lambda.zip $(BINARY_NAME)

t-apply:
	cd infrastructure/terraform ; terraform init ; terraform apply

t-destroy:
	cd infrastructure/terraform ; terraform destroy

deploy: clean build-lambda zip t-apply