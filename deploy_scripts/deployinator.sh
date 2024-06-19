#!/bin/bash

echo "Deploying to deployinator"

cd /var/www/deployinator/
git pull origin main
go mod tidy
go build -o deployinator
systemctl restart deployinator