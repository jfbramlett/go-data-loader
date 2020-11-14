db_url=jdbc\:mysql\://mariadb:3306/asapp
db_user=asapp
db_pwd=asapp
docker_opts='--network=asapp'

clean:
	rm -rf bin
	rm -rf vendor

flyway_migrate:
	docker run $(docker_opts) --rm -v $(PWD)/db/migrations\:/flyway/sql boxfuse/flyway -url=$(db_url) -user=$(db_user) -password=$(db_pwd) migrate


vendor:
	go mod vendor

build: clean vendor
	go build -o bin/server ./cmd/...

lint:
	docker run --rm -t --entrypoint=linter -v `pwd`:$(GO_PROJECT_PATH) -w $(GO_PROJECT_PATH) golang:dev-latest

build-docker:
	cp ~/.ssh/id_rsa .
	docker build -t go-template-server:latest .
	rm id_rsa

test:
	go test -cover ./pkg/...