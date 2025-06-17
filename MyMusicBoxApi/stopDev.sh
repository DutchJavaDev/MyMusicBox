#!/bin/bash

# Find and kill all the process
# Kinda overkill....
pkill -f "/tmp/go-build2246494624/b001/exe/main -port=8081 -devurl -sourceFolder=music_dev -outputExtension=opus"
pkill -f "go run main.go -port=8081 -devurl -sourceFolder=music_dev -outputExtension=opus"
pkill -f "/home/admin/.cache/go-build/cd/"
pkill -f "/b001/exe/main -port=8081 -devurl -sourceFolder=music_dev"