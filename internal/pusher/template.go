package pusher

var template = `#!/bin/sh

set -eux

docker tag $1 ${deployer}localhost:5000/$1

ssh -NT -L 5000:127.0.0.1:5000 core@10.0.122.183 &
ssh_pid=$!

docker push ${deployer}localhost:5000/$1

kill $ssh_pid
`
