net:
	docker network create auth-net

build.app:
	docker build -f Dockerfile.app -t auth .

build.db:
	docker build -f Dockerfile.db -t auth-postgres .

run.app:
	docker run --rm --name auth-container \
	--network auth-net \
	-e JWT_SECRET="jwt_secret_example" \
	-e host_db="postgres" \
	-e port_db="5432" \
	-e user_db="postgres" \
	-e password_db="postgres" \
	-e dbname_db="postgres" \
	-e sslmode_db="disable" \
	-d -p 8022:8020 auth

run.db.innet:
	docker run --rm -name auth-postgres-container \
	--network auth-net \
	-d \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_PASSWORD=postgres \
	-e POSTGRES_DB=postgres \
	-v $(pwd)/pgdata:/var/lib/postgresql/data \
	-v $(pwd)/postgres:/docker-entrypoint-initdb.d \
	-p 3333:5432 \
	auth-postgres

run.db.outnet:
	docker run --name auth-postgres-container \
    -d \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=postgres \
    -e POSTGRES_DB=postgres \
    -v $(pwd)/pgdata:/var/lib/postgresql/data \
    -v $(pwd)/postgres:/docker-entrypoint-initdb.d \
    -p 3333:5432 \
    auth-postgres

del.db:
	docker stop auth-postgres-container
	docker rm auth-postgres-container
	rm -rf $(pwd)/pgdata