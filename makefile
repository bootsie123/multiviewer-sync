EXECUTABLE=multiviewer-sync

build: build-linux build-windows build-darwin

build-linux:
	GOARCH=amd64 GOOS=linux go build -o dist/${EXECUTABLE}-linux-amd64 .
	GOARCH=arm64 GOOS=linux go build -o dist/${EXECUTABLE}-linux-arm64 .

build-windows:
	GOARCH=amd64 GOOS=windows go build -o dist/${EXECUTABLE}-windows-amd64.exe .
	GOARCH=arm64 GOOS=windows go build -o dist/${EXECUTABLE}-windows-arm64.exe .

build-darwin:
	GOARCH=amd64 GOOS=darwin go build -o dist/${EXECUTABLE}-mac-amd64 .
	GOARCH=arm64 GOOS=darwin go build -o dist/${EXECUTABLE}-mac-arm64 .

clean:
	go clean
	rm -rf dist