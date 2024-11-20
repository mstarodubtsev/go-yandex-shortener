#!/usr/bin/env bash

# go vet -vet tool=/Users/smv/work/study/go-advanced/autotests/statictest-darwin-arm64 ./...

rm shortener
go build -o shortener ./cmd/shortener/*.go
#/Users/smv/work/study/go-advanced/autotests/shortenertest-darwin-arm64 -test.v -test.run=^TestIteration1$ -binary-path=./shortener
#/Users/smv/work/study/go-advanced/autotests/shortenertest-darwin-arm64 -test.v -test.run=^TestIteration2$ -binary-path=./shortener -source-path=.
#/Users/smv/work/study/go-advanced/autotests/shortenertest-darwin-arm64 -test.v -test.run=^TestIteration3$ -binary-path=./shortener -source-path=.
/Users/smv/work/study/go-advanced/autotests/shortenertest-darwin-arm64 -test.v -test.run=^TestIteration4$ -binary-path=./shortener -source-path=.
