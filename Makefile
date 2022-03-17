.PHONY: postgres create_db drop_db migrate_up migrate_down sqlc
postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
create_db:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
drop_db:
	docker exec -it postgres12 dropdb simple_bank
link_db:
	docker exec -it postgres12 psql -U root simple_bank
migrate_init_db:
	migrate create -ext sql -dir $(GOPATH)/src/Bank/db/migration -seq init_schema
migrate_move_db:
	migrate -path $(GOPATH)/src/Bank/db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" --verbose up
migrate_up:
	migrate -path $(GOPATH)/src/Bank/db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrate_down:
	migrate -path $(GOPATH)/src/Bank/db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
migrate_version:
	migrate -path $(GOPATH)/src/Bank/db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" version
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
