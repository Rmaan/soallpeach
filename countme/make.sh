set -eu

command=$1
shift;

rate=500
duration=10
container_name=countme
image_name=countme

case $command in
run-docker)
  docker rm -fv $container_name || true
  docker build -t $image_name .
  echo $((1 + RANDOM % 10)) > payload.txt
  container_id=$(docker run -d -p "8080:80" --name $container_name $image_name)
  timeout 60 bash -c 'while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' localhost:8080/count)" != "200" ]]; do sleep 1; done' || false
  vegeta -cpus 1 attack -rate $rate -duration=${duration}s -targets target.list | vegeta report -type=json | jq '.' > metrics.json
  echo $(($(jq '.latencies["99th"]' metrics.json) / 1000000))
  ;;
*)
  echo Unknown command: $command
  ;;
esac
