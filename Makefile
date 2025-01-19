swag:
	swag fmt && swag init --parseInternal --parseDependency --parseDepth=1

prepare-db:
	docker run --name postgres-go-func -d -p 5432:5432 -e POSTGRES_USER=test -e POSTGRES_PASSWORD=password -e POSTGRES_DB=func  postgres:alpine

kill-db:
	docker rm -f postgres-go-func
