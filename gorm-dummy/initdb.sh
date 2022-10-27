#!/bin/sh

docker compose up -d

echo "Container up" 

sleep 2

USER="root"
DB="defaultdb"
FILENAME="./scripts/entity.sql"

echo ${FILENAME}

docker exec dummy_gin_db psql --username "${USER}" --dbname "${DB}" -f "${FILENAME}"
