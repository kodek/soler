#!/bin/bash
go build -o ./server_main server/server_main.go && \
	docker build -t soler/server -f Dockerfile .
