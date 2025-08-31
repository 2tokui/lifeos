build-linux:
	GOOS=linux GOARCH=amd64 go build -v -o lifeos ./cmd/lifeos/main.go

build-windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -v -o lifeos.exe ./cmd/lifeos/main.go

build: build-linux build-windows

run: build
	@./lifeos

