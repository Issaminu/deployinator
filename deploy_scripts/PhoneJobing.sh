#!/bin/bash

echo "Deploying to PhoneJobing"

cd /var/www/PhoneJobing/
git pull origin main
composer install --no-interaction
bun install
php artisan optimize
chmod -R 777 storage
chmod -R 777 bootstrap/cache
systemctl restart php-fpm
systemctl restart mysqld