.PHONY : dev prod compile

dev:
	echo "Compiling..."
	go build . 

prod:
	echo "Compiling..."
	go build -ldflags="-s -w" .

compile:
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o bin/peek-freebsd-amd64 .
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/peek-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/peek-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/peek-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/peek-darwin-m1 .