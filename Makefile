GO_VERSION := 1.24.6
APP_NAME="fritz"

up: ## Database migration up
	@go run cmd/migrations/main.go up

down: ## Database migration down
	@go run cmd/migrations/main.go down

reset: ## Database migration down
	@go run cmd/reset/main.go

reset-project: reset up

# Testing
test:
	go test ./... -coverprofile=coverage.out

coverage:
	go tool cover -func coverage.out \
	| grep "total:" | awk '{print ((int($$3) > 80) != 1) }'

report:
	go tool cover -html=coverage.out -o coverage.html

check-format:
	test -z $$(go fmt ./...)