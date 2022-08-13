.PHONY: deps test docker bin/redis-exporter

deps:
	go mod tidy

test:
	go test ./...

docker:
	docker build -t jace-ys/redis-exporter:v0.0.0 .

bin/redis-exporter:
	go build -o $@ main.go