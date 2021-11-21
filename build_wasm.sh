#!/bin/sh

set -ex

cp $(go env GOROOT)/misc/wasm/wasm_exec.js docs
GOOS=js GOARCH=wasm go build -o docs/lib.wasm wasm/main.go
