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
	docker run -d -ti --network host -v /var/www/diky2020/cdn.diky2020.cz:/upload -v ~/diky2020-backend/config.yml:/app/config.yml --name diky2020-backend diky2020-backend:1.0

dockerstop:
	docker stop diky2020-backend

dockerrm:
	docker container rm diky2020-backend
