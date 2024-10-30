#!/usr/bin/env bash

go build -o shortener ./cmd/shortener/main.go
/Users/smv/work/study/go-advanced/autotests/shortenertest-darwin-arm64 -test.v -test.run=^TestIteration1$ -binary-path=./shortener
