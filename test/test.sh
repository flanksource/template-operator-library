#!/bin/bash

set -ex

export PLATFORM_CLI_VERSION=v0.30.0
export PLATFORM_CLI="./karina -c test/config.yaml"
export KUBECONFIG=~/.kube/config
export DOCKER_API_VERSION=1.39

make

if which karina 2>&1 > /dev/null; then
  PLATFORM_CLI="karina -c test/config.yaml"
else
  if [[ "$OSTYPE" == "linux-gnu" ]]; then
    wget -q https://github.com/flanksource/karina/releases/download/$PLATFORM_CLI_VERSION/karina
    chmod +x karina
  elif [[ "$OSTYPE" == "darwin"* ]]; then
    wget -q https://github.com/flanksource/karina/releases/download/$PLATFORM_CLI_VERSION/karina_osx
    cp karina_osx karina
    chmod +x karina
  else
    echo "OS $OSTYPE not supported"
    exit 1
  fi
fi

mkdir -p .bin

kind delete cluster --name kind-kind
$PLATFORM_CLI ca generate --name root-ca --cert-path .certs/root-ca.crt --private-key-path .certs/root-ca.key --password foobar  --expiry 1
$PLATFORM_CLI ca generate --name ingress-ca --cert-path .certs/ingress-ca.crt --private-key-path .certs/ingress-ca.key --password foobar  --expiry 1
$PLATFORM_CLI provision kind-cluster -v 5 --trace

$PLATFORM_CLI deploy crds
$PLATFORM_CLI deploy calico
$PLATFORM_CLI deploy base
$PLATFORM_CLI deploy minio
$PLATFORM_CLI deploy stubs
$PLATFORM_CLI deploy template-operator
$PLATFORM_CLI deploy redis-operator
$PLATFORM_CLI deploy eck

cd test
go test ./... -v
