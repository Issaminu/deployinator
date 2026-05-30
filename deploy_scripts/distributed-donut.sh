#!/bin/bash
set -euo pipefail

echo "Deploying distributed-donut"

cd /var/www/distributed-donut/
trap 'rm -f /var/www/distributed-donut/donut-server.new' EXIT

as_repo_user() {
	if [ "$(id -u)" -eq 0 ]; then
		runuser -u issam -- "$@"
	else
		"$@"
	fi
}

if [ -n "$(as_repo_user git status --porcelain)" ]; then
	echo "Refusing to deploy: /var/www/distributed-donut has local changes"
	exit 1
fi

as_repo_user git fetch origin main
as_repo_user git reset --hard origin/main
# as_repo_user go test ./...
as_repo_user go build -o donut-server.new ./cmd/donut-server
as_repo_user mv donut-server.new donut-server
systemctl restart donut
systemctl is-active --quiet donut
curl -fsS --max-time 2 http://127.0.0.1:8080/ >/dev/null

echo "distributed-donut deployed"
