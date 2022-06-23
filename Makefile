server : $(wildcard cmd/server//*.go) $(wildcard internal/firewall/*.go)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o server ./cmd/server/

client : cmd/client/main.go
	go build -o client ./cmd/client/
