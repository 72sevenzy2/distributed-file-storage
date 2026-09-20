#!/bin/bash

cd .. || exit 1
go test pathTransform_test.go storage.go -v
