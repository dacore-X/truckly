createdb:
	docker exec -it postgres15 createdb --username=root --owner=root truckly

dropdb:
	docker exec -it postgres15 dropdb truckly

psql:
	docker exec -it postgres15 psql -h localhost -p 5432 -U root -W

migrateup:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/truckly?sslmode=disable" -verbose up

migratedown:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/truckly?sslmode=disable" -verbose down

run:
	go run cmd/app/main.go

.PHONY: createdb dropdb psql migrateup migratedown run