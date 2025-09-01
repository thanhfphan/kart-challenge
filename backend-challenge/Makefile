start:
	@./scripts/bin.sh start
.PHONY: start

tests:
	@./scripts/bin.sh tests
.PHONY: tests

up:
	@./scripts/bin.sh up
.PHONY: up

down:
	@./scripts/bin.sh down
.PHONY: down

migrate:
	@./scripts/bin.sh migrate
.PHONY: migrate

preprocess:
	@echo "Running promo code preprocessing pipeline..."
	@go run cmd/preprocess/main.go -output_dir=data
.PHONY: preprocess

help:
	@./scripts/bin.sh help
.PHONY: help
