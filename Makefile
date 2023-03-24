createdb:
	docker exec -it postgres15 createdb --username=root --owner=root truckly

dropdb:
	docker exec -it postgres15 dropdb truckly

migrateup:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/truckly?sslmode=disable" -verbose up

migratedown:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/truckly?sslmode=disable" -verbose down

.PHONY: createdb dropdb migrateup migratedown