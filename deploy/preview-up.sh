#!/usr/bin/env bash
# Create or update the preview backend for one pull request.
#
# Runs ON the droplet, fed over stdin so the droplet needs no checkout:
#   ssh host bash -s -- <pr> <backend-image> <migrate-image> < deploy/preview-up.sh
#
# Secrets are NOT arguments: they are read from /opt/<project>/preview.env, so
# they never appear in the droplet's process list or in CI logs.
set -euo pipefail

PR_NUMBER=${1:?usage: preview-up.sh <pr-number> <backend-image> <migrate-image>}
BACKEND_IMAGE=${2:?missing backend image}
MIGRATE_IMAGE=${3:?missing migrate image}

PROJECT=inspirate-consulting
ENV_FILE="/opt/${PROJECT}/preview.env"
PREVIEW_DOMAIN_FILE="/opt/${PROJECT}/preview-domain"

[ -r "$ENV_FILE" ] || { echo "missing $ENV_FILE" >&2; exit 1; }
[ -r "$PREVIEW_DOMAIN_FILE" ] || { echo "missing $PREVIEW_DOMAIN_FILE" >&2; exit 1; }

PREVIEW_DOMAIN=$(cat "$PREVIEW_DOMAIN_FILE")
DB_NAME="pr_${PR_NUMBER}"
CONTAINER="${PROJECT}-pr-${PR_NUMBER}"
HOST="pr-${PR_NUMBER}.${PREVIEW_DOMAIN}"

# Read the Postgres superuser password out of the env file rather than taking
# it as an argument.
PGPASSWORD=$(grep -E '^DB_PASSWORD=' "$ENV_FILE" | head -1 | cut -d= -f2-)
[ -n "$PGPASSWORD" ] || { echo "DB_PASSWORD is empty in $ENV_FILE" >&2; exit 1; }

echo "==> preview for PR #${PR_NUMBER} at https://${HOST}"

# 1. Database. Idempotent: a re-push to the same PR reuses the existing one so
#    that data entered during review survives.
if docker exec -e PGPASSWORD="$PGPASSWORD" preview-db \
     psql -U postgres -tAc "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}'" | grep -q 1; then
  echo "    database ${DB_NAME} already exists"
else
  echo "    creating database ${DB_NAME}"
  docker exec -e PGPASSWORD="$PGPASSWORD" preview-db createdb -U postgres "${DB_NAME}"
fi

# 2. Migrations, before the new code starts, so a bad migration fails the
#    preview instead of leaving the app against a schema it cannot read.
echo "    applying migrations"
docker pull --quiet "$MIGRATE_IMAGE"
docker run --rm --network edge \
  -e DATABASE_URL="postgres://postgres:${PGPASSWORD}@preview-db:5432/${DB_NAME}?sslmode=disable" \
  "$MIGRATE_IMAGE"

# 3. Backend container. Replacing it is how an update happens; Caddy notices
#    the new container through the Docker socket and reroutes automatically.
echo "    starting ${CONTAINER}"
docker pull --quiet "$BACKEND_IMAGE"
docker rm -f "$CONTAINER" >/dev/null 2>&1 || true
docker run -d \
  --name "$CONTAINER" \
  --network edge \
  --restart unless-stopped \
  --env-file "$ENV_FILE" \
  -e DB_NAME="$DB_NAME" \
  --label "caddy=${HOST}" \
  --label "caddy.reverse_proxy={{upstreams 8080}}" \
  "$BACKEND_IMAGE" >/dev/null

echo "==> https://${HOST}"
