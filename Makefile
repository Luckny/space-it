include ./app.env

postgres:
	docker run --name postgresql16 -p 5432:5432 -e POSTGRES_USER=${DB_ADMIN_USER} -e POSTGRES_PASSWORD=${DB_ADMIN_PASSWORD} -d postgres:16-alpine

createdb:
	docker exec -it postgresql16 createdb --username=${DB_ADMIN_USER} --owner=${DB_ADMIN_USER} ${DB_NAME}

createapiuser:
	docker exec -it postgresql16 psql -U ${DB_ADMIN_USER} -c "CREATE USER "${DB_USER}" WITH ENCRYPTED PASSWORD '"${DB_PASSWORD}"';"

dropapiuser:
	docker exec -it postgresql16 psql -U ${DB_ADMIN_USER} -c "DROP USER IF EXISTS "${DB_USER}";"

dropdb:
	docker exec -it postgresql16 dropdb ${DB_NAME}

migrateup:
	migrate -path db/migration -database "postgres://"${DB_ADMIN_USER}":"${DB_ADMIN_PASSWORD}"@"${DB_HOST}":"${DB_PORT}"/"${DB_NAME}"?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://"${DB_ADMIN_USER}":"${DB_ADMIN_PASSWORD}"@"${DB_HOST}":"${DB_PORT}"/"${DB_NAME}"?sslmode=disable" -verbose down


.PHONY:
	postgres createdb createapiuser dropapiuser dropdb migrateup migratedown
