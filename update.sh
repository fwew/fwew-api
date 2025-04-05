#!/bin/bash
git pull
rm ./fwew-api
go build ./...
source .env
./fwew-api
