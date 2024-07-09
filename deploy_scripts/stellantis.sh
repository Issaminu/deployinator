#!/bin/bash

echo "Deploying Stellantis VENG Project Manager"

cd /var/www/stellantis-project-manager/
git pull origin main
systemctl restart postgresql
npm install
npx prisma db push --accept-data-loss
npx prisma generate
npx @digitak/esrun seed/seed.ts
npx next build
npx next start --port 7887
systemctl restart stellantis