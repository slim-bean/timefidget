
build-markup:
	go build -o cmd/markup/markup ./cmd/markup/main.go

build-notifier-arm:
	GOOS=linux GOARCH=arm GOARM=7 go build -o cmd/notifier/notifier ./cmd/notifier/main.go