DB_HOST := $(shell yq '.migrate_db_host' config.yaml )
DB_PORT := $(shell yq '.database.port' config.yaml)
DB_USER := $(shell yq '.database.username' config.yaml)
DB_PASSWORD := $(shell yq '.database.password' config.yaml)
DB_NAME := $(shell yq '.database.dbname' config.yaml)
DB_SSLMODE := $(shell yq '.database.sslmode' config.yaml)
SERVER_HOST := $(shell yq '.server.host' config.yaml)
SERVER_PORT := $(shell yq '.server.port' config.yaml)

MIGRATIONS_FILE := "./migrations"
DB_URL := "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"

migrate-up: 
	goose -dir $(MIGRATIONS_FILE) postgres $(DB_URL) up

migrate-down:
	goose -dir $(MIGRATIONS_FILE) postgres $(DB_URL) down

create-migration:
	goose create $(name) sql -dir ./migrations

start:
	sudo docker-compose up --build -d
	
stop: 
	sudo docker-compose stop
