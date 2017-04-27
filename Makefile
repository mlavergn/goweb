###############################################
#
# Makefile
#
###############################################

FLAGS  :=

all: fmt build

fmt:
	go fmt *.go

build:
	go build

test:
	go test

run:
	go run

clean:
	go clean
