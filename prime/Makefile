prime: prime.go
	go build -o prime prime.go

run: prime
	time ./prime input.txt > output.txt

run-docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o prime-docker prime.go
	docker build -t prime .
	rm -rf data/
	mkdir data/
	cp input.txt data/in.txt
	time docker run --rm -v "$(realpath data):/data" prime /data/in.txt /data/out.txt > output.txt
	if [[ ! -s output.txt && -s data/out.txt ]]; then mv data/out.txt output.txt; fi

diff:
	diff -w output.txt expected.txt

