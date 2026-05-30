#!/bin/bash
set -euo pipefail

echo "Deploying distributed-donut"

cd /var/www/distributed-donut/

as_repo_user() {
	if [ "$(id -u)" -eq 0 ]; then
		runuser -u issam -- "$@"
	else
		"$@"
	fi
}

as_repo_user git pull --ff-only origin main
// as_repo_user go test ./...
as_repo_user go build -o donut-server ./cmd/donut-server
systemctl restart donut
systemctl is-active --quiet donut
curl -fsS --max-time 2 http://127.0.0.1:8080/ >/dev/null

echo "distributed-donut deployed"
