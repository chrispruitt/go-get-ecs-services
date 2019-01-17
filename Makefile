VERSION=`go run main.go -version`
default:
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o build/get-service-versions.$(VERSION).windows-amd64.exe
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/get-service-versions.$(VERSION).linux-amd64
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o build/get-service-versions.$(VERSION).darwin-amd64

build: default

tag:
	git tag $(VERSION)
