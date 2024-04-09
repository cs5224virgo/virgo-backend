POSTGRES_URL = 'postgres://virgo:virgo@127.0.0.1:5432/virgo?sslmode=disable'
SHELL = /bin/bash
CURRENT_MIGRATION_VER = $(shell migrate -database ${POSTGRES_URL} -path db/migrations version 2>&1| cut -f 1 -d ' ')
PREVIOUS_MIGRATION_VER = $(shell expr ${CURRENT_MIGRATION_VER} - 1)
NEXT_MIGRATION_VER = $(shell expr ${CURRENT_MIGRATION_VER} + 1)

build:
	go build

clean:
	rm -r db/generated
	rm go.sum

generate:
	go generate

migrate-up:
	migrate -database ${POSTGRES_URL} -path db/migrations up

migrate-down:
	migrate -database ${POSTGRES_URL} -path db/migrations down

migrate-force-previous:
	migrate -database ${POSTGRES_URL} -path db/migrations force ${PREVIOUS_MIGRATION_VER}

migrate-force-current:
	migrate -database ${POSTGRES_URL} -path db/migrations force ${CURRENT_MIGRATION_VER}

migrate-force-next:
	migrate -database ${POSTGRES_URL} -path db/migrations force ${NEXT_MIGRATION_VER}

migrate-version:
	migrate -database ${POSTGRES_URL} -path db/migrations version

migrate-drop:
	migrate -database ${POSTGRES_URL} -path db/migrations drop
