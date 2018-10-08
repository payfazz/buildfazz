package template

var BuilderScript = map[string]string{
	// go builder
	"docker-builder-go": `FROM golang:alpine

RUN set -eux \
 && wget -qO /usr/local/bin/docker_pid1 https://github.com/win-t/docker_pid1/releases/download/v3.0.0/docker_pid1 \
 && chmod 755 /usr/local/bin/docker_pid1

RUN set -eux \
 && apk -U add gcc git musl-dev rsync

RUN set -eux \
 && adduser -D user \
 && mkdir -p ~user/go \
 && chown -R user:user ~user/go

USER user

ENV GOPATH=/home/user/go

ENV PATH=/home/user/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

WORKDIR /home/user/go

ENTRYPOINT ["docker_pid1"]`,

	// alpine builder
	"docker-alpine-builder": `FROM alpine:latest

RUN set -eux \
 && wget -qO /usr/local/bin/docker_pid1 https://github.com/win-t/docker_pid1/releases/download/v3.0.0/docker_pid1 \
 && chmod 755 /usr/local/bin/docker_pid1 \
 && wget -qO /usr/local/bin/swuser https://github.com/win-t/switch-user/releases/download/v1.0.0/swuser \
 && chmod 755 /usr/local/bin/swuser

RUN set -eux \
 && mkdir /logs \
 && ln -sf /dev/stdout /logs/out.txt \
 && ln -sf /dev/stderr /logs/err.txt

RUN set -eux \
 && apk add -U ca-certificates

RUN set -eux \
 && adduser -D app

ADD output /app
ADD scripts/start-app /
ENTRYPOINT ["/start-app"]`,

	// build
	"build": `#!/bin/sh

set -eu

cd "$(dirname "$0")"

docker_tag=${DOCKER_TAG:-}

if [ -z "$docker_tag" ]; then
  echo "DOCKER_TAG is not set" >&2
  exit 1
fi

./builder start
./sync-src
./builder exec sh -ceu 'cd ~ && ./build "$@"' - "$@"
./sync-output
./builder rm -f

pull_opt=--pull
[ "${NO_PULL:-}" = y ] && pull_opt=
cache_opt=
[ "${NO_CACHE:-}" = y ] && cache_opt=--no-cache
tar -c \
  output \
  scripts/start-app \
  dockerfiles/project \
| docker build $pull_opt $cache_opt \
  -t "$docker_tag" \
  -f dockerfiles/project -`,

	// builder
	"builder": `#!/usr/bin/env docker.sh

name=@project-name$

image=$name

opts="
  -v '$name-data:/home/user'
"

command_rmvol() {
  docker volume rm "$name-data" > /dev/null
}

pre_pull() {
  panic "Image $image can't be pulled, please build"
}

command_build() (
  pull_opt=--pull
  [ "${NO_PULL:-}" = y ] && pull_opt=

  cache_opt=
  [ "${NO_CACHE:-}" = y ] && cache_opt=--no-cache

  cd "$dir/dockerfiles"
  tar -c builder | docker build $pull_opt $cache_opt \
    -t "$image" \
    -f builder -
)

pre_start() {
  if [ "$1" = run ] && ! exists image "$image"; then
    command_build || panic "Failed to build builder image"
  fi
}`,

	// sync output
	"sync-output": `#!/bin/sh

set -eu

cd "$(dirname "$0")"

alias docker_sync="rsync -e 'docker exec -i' --blocking-io"

mkdir -p output

docker_sync \
  -avr --delete \
  "$(./builder name):/home/user/output/" \
  output`,

	// sync src
	"sync-src": `#!/usr/bin/env bash

#!/bin/sh

set -eu

cd "$(dirname "$0")"

alias docker_sync="rsync -e 'docker exec -i' --blocking-io"
./builder exec sh -ceu 'cd ~/go && mkdir -p "src/github.com/payfazz/@project-name$"'

docker_sync \
  -a --delete \
  scripts/install-deps \
  scripts/build \
  "$(./builder name):/home/user"

docker_sync \
  --exclude="/.git/" \
  --exclude="/docker/" \
  -a --delete \
  ../../ \
  "$(./builder name):/home/user/go/src/github.com/payfazz/@project-name$"
`,

	//script build go
	"scripts-build-go": `#!/bin/sh

set -eu

cd "$(dirname "$0")"

lock_dir="$PWD/$(basename "$0").lock"

if ! mkdir "$lock_dir" >/dev/null 2>&1; then
  echo "cannot aquire lock" >&2
  exit 1
fi

trap "rmdir \"$lock_dir\"" EXIT

get_flag=-u

for arg; do
  case $arg in
  --no-update) get_flag= ;;
  esac

done
(
  set -x
  cd ~
  outdir=$PWD/output
  mkdir -p "$outdir"
  install-deps-go "$@"
  cd "$(go env GOPATH)/src/github.com/payfazz/@project-name$"
  GOBIN=$outdir go install ./...
)`,

	//install dep
	"scripts-install-dep-go": `#!/bin/sh

set -eu

cd "$(dirname "$0")"

lock_dir="$PWD/$(basename "$0").lock"

if ! mkdir "$lock_dir" >/dev/null 2>&1; then
  echo "cannot aquire lock" >&2
  exit 1
fi

trap "rmdir \"$lock_dir\"" EXIT

get_flag=-u

for arg; do
  case $arg in
  --no-update) get_flag= ;;
  esac

done
(
  cd "$(go env GOPATH)/src/github.com/payfazz/@project-name$"
  pkgname=$(go list .)
  getpkgs=$(
    go list -f '{{range .Deps}}{{.|printf "%s\n"}}{{end}}' ./... \
    | grep -v "^$pkgname\$" | grep -v "^$pkgname/" \
    | sort | uniq
  )
  go get $get_flag $getpkgs
)`,

	// start app
	"scripts-start-app-go": `#!/bin/sh

set -eu

if [ "$(id -u)" = "0" ]; then
    chown app:app -R /logs
    export HOME=$(getent passwd app | cut -d: -f6)
    export USER=app
    export LOGNAME=app
    exec 1>>/logs/out.txt
    exec 2>>/logs/err.txt
    exec docker_pid1 swuser $(id -u app),$(id -g app) "$0"
    exit 1
else
    cd /app
    exec ./@project-name$
fi`,
}
