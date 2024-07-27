postgres:
	sudo docker run --name sfss-storage -e POSTGRES_PASSWORD=12345678 -d postgres

createdb:
	sudo docker exec -it sfss-storage createdb --username postgres --owner=postgres sfss-storage

dropdb:
	sudo docker exec -it sfss-storage dropdb sfss-storage

migrateup:
	migrate -path db/migrations -database "postgresql://postgres:12345678@172.17.0.2:5432/sfss-storage?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:12345678@172.17.0.2:5432/sfss-storage?sslmode=disable" -verbose down

.PHONY:
	postgres createdb dropdb migrateup migratedown
