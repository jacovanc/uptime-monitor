build:
	go build -o bin/uptime-monitor.exe cmd/main.go

run:
	CGO_ENABLED=1 go run cmd/main.go

test:
	go test ./...
