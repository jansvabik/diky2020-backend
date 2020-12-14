default:

up:
	go run cmd/backend/main.go

main:
	go build cmd/backend/main.go

run: main
	./main
