prime: prime.go
	go build -o prime prime.go

run: prime
	time ./prime input.txt > output.txt

run-docker:
	docker build -t prime .
	rm -rf data/
	mkdir data/
	cp input.txt data/in.txt
	time docker run --rm -v "$(realpath data):/data" --entrypoint= prime sh -c '/prime /data/in.txt > /data/out.txt'
	mv data/out.txt output.txt

diff:
	diff -w output.txt expected.txt

