#!/usr/bin/env bash

set -euo pipefail

pkg="ghe-get-all-owners"
platforms=("windows/amd64" "windows/386" "darwin/amd64")

for platform in "${platforms[@]}"; do
    # shellcheck disable=SC2206
    platform_split=(${platform//\// })

    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    output_name="${pkg}"

    if [ "${GOOS}" = "windows" ]; then
        output_name+="-${GOARCH}.exe"
    fi

    [ ! -d "./build/${GOOS}" ] && mkdir -p "./build/${GOOS}"

    env GOOS="${GOOS}" GOARCH="${GOARCH}" go build -o "./build/${GOOS}/${output_name}" . || {
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    }
done
