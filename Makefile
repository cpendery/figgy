test:
	go test ./... -coverprofile=coverage.out
	go tool cover --html=coverage.out

snap:
	UPDATE_SNAPSHOTS=true go test ./...

build:
	go build -o fig ./cmd
