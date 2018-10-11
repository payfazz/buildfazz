package pusher

var template = `#!/bin/sh

set -eux

docker tag $1 ${deployer}${server}/$1

ssh -NT -L 5000:${server} ${ssh} &
ssh_pid=$!

docker push ${deployer}${server}/$1

kill $ssh_pid
`
