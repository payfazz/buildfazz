package pusher

var template = `#!/bin/sh

set -u

panic() {
  panic_code=$1
  shift
  echo "$@" >&2
  exit $panic_code
}

docker tag $1 ${deployer}${server}/$1 || panic $? "cannot tag"

ssh -NT -L 5000:${server} ${ssh} &
ssh_pid=$!

sleep 1 

exit_code=0
docker push ${deployer}${server}/$1 || exit_code=$?

kill $ssh_pid

exit $exit_code
`
