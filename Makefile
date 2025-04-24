## Filename Makefile
include .envrc

# .PHONY: run/tests
# run/tests: vet
# 	go test -v ./...

.PHONY: fmt
fmt: 
	go fmt ./...

.PHONY: vet
vet: fmt
	go vet ./...

.PHONY: run
run: 
	go run ./cmd/web -addr=${ADDRESS} -dsn=${TEST_1_DB_DSN}

.PHONY: db/psql
db/psql:
	psql ${TEST_1_DB_DSN}


## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}


# ## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${TEST_1_DB_DSN} up

# # db/migrations/down-1: undo the last migration
.PHONY: db/migrations/down-1
db/migrations/down-1:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${TEST_1_DB_DSN} down 1

