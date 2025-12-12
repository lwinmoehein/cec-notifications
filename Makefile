.PHONY: build clean deploy test deps local-invoke

# Build the Lambda binary
build:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go

# Install dependencies
deps:
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	rm -f bootstrap
	rm -rf .aws-sam

# Build and deploy
deploy: build
	sam deploy --guided

# Deploy without prompts (requires samconfig.toml)
deploy-fast: build
	sam deploy

# Run tests
test:
	go test -v ./...

# Test locally with SAM
local-invoke: build
	sam local invoke NotificationFunction --event events/sample-sqs-event.json

# Start SAM local API
local-api: build
	sam local start-api

# Tail CloudWatch logs
logs:
	sam logs -n NotificationFunction --tail

# Validate SAM template
validate:
	sam validate

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run
