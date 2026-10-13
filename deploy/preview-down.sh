#!/usr/bin/env bash
# Tear down the preview backend for one pull request.
#
#   ssh host bash -s -- <pr> < deploy/preview-down.sh
#
# Best-effort by design: a half-cleaned preview must not fail the workflow, or
# leftovers would accumulate with no way to clear them.
set -uo pipefail

PR_NUMBER=${1:?usage: preview-down.sh <pr-number>}

PROJECT=inspirate-consulting
ENV_FILE="/opt/${PROJECT}/preview.env"
DB_NAME="pr_${PR_NUMBER}"
CONTAINER="${PROJECT}-pr-${PR_NUMBER}"

echo "==> tearing down preview for PR #${PR_NUMBER}"

docker rm -f "$CONTAINER" >/dev/null 2>&1 \
  && echo "    removed ${CONTAINER}" \
  || echo "    no container ${CONTAINER}"

if [ -r "$ENV_FILE" ]; then
  PGPASSWORD=$(grep -E '^DB_PASSWORD=' "$ENV_FILE" | head -1 | cut -d= -f2-)
  if [ -n "$PGPASSWORD" ]; then
    # FORCE disconnects any session still holding the database open, which a
    # just-killed container can briefly do.
    docker exec -e PGPASSWORD="$PGPASSWORD" preview-db \
      psql -U postgres -tAc "DROP DATABASE IF EXISTS ${DB_NAME} WITH (FORCE)" >/dev/null 2>&1 \
      && echo "    dropped ${DB_NAME}" \
      || echo "    could not drop ${DB_NAME} (may not exist)"
  fi
fi

echo "==> done"
