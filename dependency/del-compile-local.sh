#!/usr/bin/env bash
set -euo pipefail

if [[ "${1}" == "bionic" ]]; then
    docker build -t compilation -f actions/compile/bionic.Dockerfile actions/compile/
    docker run --rm -v /tmp/dep-output/:"${PWD}" compilation --version "2.4.54" --outputDir "${PWD}" --target "bionic"
    exit 0
fi

if [[ "${1}" == "jammy" ]]; then
    docker build -t compilation -f actions/compile/jammy.Dockerfile actions/compile/
    docker run --rm -v /tmp/dep-output/:"${PWD}" compilation --version "2.4.54" --outputDir "${PWD}" --target "jammy"
    exit 0
fi

if [[ "${1}" == "test" ]]; then
    pushd "test"
        docker build -t test .
        docker run --rm -v /tmp/dep-output/:/tarball_path test --version "2.4.54" # docker run -it --rm -v /tmp/dep-output/:/tarball_path --entrypoint bash test
    popd
    exit 0
fi
