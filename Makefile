test:
	go clean -testcache
	go test ./... -v
run:
	go run "cmd/server/main.go"