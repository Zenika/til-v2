#!/bin/bash

echo "Starting to build til-aio image..."

echo "Building server"
cd ./server
docker build -t til-backend:latest .
echo "Building client"
cd ../ui
docker build -t til-frontend:latest .
echo "Building AIO based on previous build"
cd ../
docker build -t til-aio:latest .

echo "Build done; you should now be able to run til-aio:latest."