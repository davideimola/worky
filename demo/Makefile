.PHONY: run status check reset build tidy

tidy:
	go mod tidy

run:
	go run . serve --open

status:
	go run . status

check:
	go run . check

reset:
	go run . reset

build:
	go build -o bin/demo .
