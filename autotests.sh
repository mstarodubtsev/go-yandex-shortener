#!/usr/bin/env bash

rm shortener
go build -o shortener ./cmd/shortener/*.go
/Users/smv/work/study/go-advanced/autotests/shortenertest-darwin-arm64 -test.v -test.run=^TestIteration1$ -binary-path=./shortener
