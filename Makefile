clean:
	go clean

build:
	.tools/build.sh default

build-linux:
	.tools/build.sh linux

test:
	go test ./...

lint:
	golangci-lint -v run

depends-up:
	docker-compose up -d --remove-orphans

depends-stop:
	docker-compose stop