#!/bin/bash

go build CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .
