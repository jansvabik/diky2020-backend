default:

up:
	go run cmd/backend/main.go

build:
	go build cmd/backend/main.go

run: build
	./main

dockerbuild:
	docker build -f Dockerfile -t diky2020-backend:1.0 .

dockerrun:
	docker run -d -ti -p 9000:80 --name diky2020-backend diky2020-backend:1.0

dockerstop:
	docker stop diky2020-backend

dockerrm:
	docker container rm diky2020-backend
