#!/bin/bash

#go run cmd/latte-compiler/main.go custom.lat
docker run -it -v "$PWD":/usr/src/comp -w /usr/src/comp gcc:11.2.0 bash compile_custom.sh
