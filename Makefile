.PHONY: postgres create_db drop_db migrate_up migrate_down migrate_up1 migrate_down1 sqlc test server mock_install mock_init
postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
create_db:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
drop_db:
	docker exec -it postgres12 dropdb simple_bank
link_db:
	docker exec -it postgres12 psql -U root simple_bank
migrate_install:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz & sudo mv migrate /usr/bin/migrate
migrate_init_db:
	migrate create -ext sql -dir $(GOPATH)/src/Bank/db/migration -seq init_schema
migrate_move_db:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" --verbose up
migrate_up:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrate_up1:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migrate_down:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
migrate_down1:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
migrate_version:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" version
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
mock_install:
	go install github.com/golang/mock/mockgen@v1.6.0 &
mock_init:
	mockgen -packge mockdb -destination db/mock/store.go github.com/0RAJA/Bank/db/sqlc Store
server:
	go run cmd/bank/bank.go
