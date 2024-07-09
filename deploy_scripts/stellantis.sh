#!/bin/bash

echo "Deploying Stellantis VENG Project Manager"

cd /var/www/stellantis-project-manager/
git reset --hard
git pull origin main
systemctl restart postgresql
npm install
npx prisma db push --accept-data-loss
npx prisma generate
npx @digitak/esrun seed/seed.ts
npx next build
systemctl restart stellantis