FROM golang:latest

RUN apt-get install -y gcc g++ make

# Docker image for running tests. This image is needed because tests use SQLite3 as in-memory database
# and that requires CGO to be enabled, which in turn requires GCC and G++ to be installed.

WORKDIR /src
ADD go.mod .
ADD go.sum .

RUN go install ./...
