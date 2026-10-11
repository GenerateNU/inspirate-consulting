#!/bin/sh
# Emits the _redirects file that Netlify uses for routing.
#
# /api is proxied to the backend so the browser only ever talks to one origin.
# That means no CORS, and no absolute API URL baked into the JS bundle -- which
# is what lets one build run unchanged on production and on any preview URL.
#
# API_PROXY_TARGET comes from the Netlify site settings (or the release
# workflow). Netlify supports per-context values, so when per-PR backends exist
# the deploy-preview context can point at its own backend without a code change.
set -eu

: "${API_PROXY_TARGET:?set API_PROXY_TARGET in Netlify site settings}"

# Order matters: the SPA fallback would otherwise swallow API calls.
cat > dist/_redirects <<REDIRECTS
/api/*  ${API_PROXY_TARGET}/:splat  200
/*      /index.html                 200
REDIRECTS

echo "wrote dist/_redirects proxying /api -> ${API_PROXY_TARGET}"
