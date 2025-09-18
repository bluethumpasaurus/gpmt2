#!/bin/bash

SOURCE="github.com/bluethumpasaurus/gpmt2/cmd/gpmt"
TARGET="build/gpmt"

echo "building binary at build/gpmt"

go build -o "${TARGET}" "${SOURCE}"
