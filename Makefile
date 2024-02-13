build:
	go build -o bin/uptime-monitor.exe cmd/uptime-monitor/main.go

run:
	go run cmd/uptime-monitor/main.go

test:
	go test ./...
