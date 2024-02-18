build:
	go build -o bin/uptime-monitor.exe cmd/main.go

run:
	go run cmd/main.go

test:
	go test ./...
