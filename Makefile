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
ifndef name
	go test
else
	go test -run $(name)
endif

run:
	go run

clean:
	go clean
