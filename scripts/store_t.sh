#!/bin/bash

cd .. || exit 1
go test storage_test.go storage.go -v
