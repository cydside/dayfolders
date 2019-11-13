#!/bin/bash

go build -o dayfolders dayfolders.go
GOOS=windows GOARCH=amd64 go build -o dayfolders_amd64.exe dayfolders.go
GOOS=windows GOARCH=386 go build -o dayfolders_386.exe dayfolders.go
