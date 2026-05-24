swag:
	swag init -g server.go \
		--dir api-server,api-server/controller/api/swagger,api-server/request-models,api-server/contracts \
		--output docs

build:
	go build ./...
