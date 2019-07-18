#! /bin/bash

# Build web and other services

cd ~/go/src/github.com/streaming_media/api
env GOOS=linux GOARCH=amd64 go build -o ../bin/api

cd ~/go/src/github.com/streaming_media/scheduler
env GOOS=linux GOARCH=amd64 go build -o ../bin/scheduler

cd ~/go/src/github.com/streaming_media/streamserver
env GOOS=linux GOARCH=amd64 go build -o ../bin/streamserver

cd ~/go/src/github.com/streaming_media/web
env GOOS=linux GOARCH=amd64 go build -o ../bin/web