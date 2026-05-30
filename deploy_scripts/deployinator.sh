#!/bin/bash
set -euo pipefail

echo "Deploying to deployinator"

cd /var/www/deployinator/
trap 'rm -f /var/www/deployinator/deployinator.new' EXIT

as_repo_user() {
	if [ "$(id -u)" -eq 0 ]; then
		runuser -u issam -- "$@"
	else
		"$@"
	fi
}

as_repo_user git pull --ff-only origin main
as_repo_user go mod tidy
as_repo_user go build -buildvcs=false -o deployinator.new
as_repo_user mv deployinator.new deployinator
systemctl restart deployinator
