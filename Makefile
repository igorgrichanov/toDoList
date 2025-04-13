run:
	docker compose up

down:
	docker compose down

swagger:
	swag fmt
	swag init -g internal/controller/http/router.go --parseInternal
