#/bin/bash

echo "Deploying to vectra-backend-ai"

cd /var/www/vectra-backend-ai/
git pull origin main
pip install -r requirements.txt
python orm.py || (rm database.db -f && python orm.py)
systemctl restart vectra-backend-ai