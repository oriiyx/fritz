APP_NAME="fritz"

up: ## Database migration up
	@go run cmd/migrations/main.go up

down: ## Database migration down
	@go run cmd/migrations/main.go down

reset: ## Database migration down
	@go run cmd/reset/main.go

reset-project: reset up
