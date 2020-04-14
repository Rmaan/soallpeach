prime: prime.go
	go build -o prime prime.go

run: prime
	time ./prime input.txt > output.txt

run-docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o prime-docker prime.go
	docker build -t prime .
	time docker run --rm -v "$(realpath input.txt):/in" prime /in > output.txt
