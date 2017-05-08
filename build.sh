#!/usr/bin/env bash
[[ -n "build" ]] && rm -rf "build"
echo "building"
go build -o build/geoip-seeker-server cmd/seeker/main.go
echo "copy files"
cp -r assets/data assets/seeker.json build
cp assets/geoip-seeker-server.service.example build/geoip-seeker-server.service
echo "build done"