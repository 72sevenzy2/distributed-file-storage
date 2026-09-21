#!/bin/bash

cd .. || exit 1
go test -run '^TestStorage$' -v
