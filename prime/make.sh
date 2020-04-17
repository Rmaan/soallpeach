set -eu

command=$1
shift;

case $command in
run)
  go build -o prime prime.go
  time ./prime input.txt >output.txt
  ;;
diff)
  diff -w output.txt expected.txt
  ;;
run-docker)
  docker build -t prime .
  rm -rf data/
  mkdir data/
  cp input.txt data/in.txt
  cid=$(docker create --rm -v "$PWD/data:/data" --entrypoint= prime sh -c '/prime /data/in.txt > /data/out.txt')
  docker start $cid
  time docker wait $cid
  mv data/out.txt output.txt
  ;;
*)
  echo Unknown command: $command
  ;;
esac
