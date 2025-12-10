APP_NAME="fritz"

up: ## Database migration up
	@go run cmd/migrations/main.go up

down: ## Database migration down
	@go run cmd/migrations/main.go down

reset-project: down up
