# Makefile

all: build lint tidy

travis-ci: build unittest lint tidy

build:
	go build ./vpcv1

unittest:
	cd vpcv1 && go test -run TestVpcV1 && cd ..

alltest:
	go test `go list ./... | grep vpcv1` -v -tags=integration -skipForMockTesting -testCount

lint:
	golangci-lint run

tidy:
	go mod tidy