# Deployment

| Piece | Runs on | Ships when |
| --- | --- | --- |
| Backend (production) | Shared DigitalOcean droplet | A GitHub Release is published |
| Frontend (production) | Netlify | A GitHub Release is published |
| Backend (PR preview) | Same droplet, one container + database per PR | Automatically, on every PR |
| Frontend (PR preview) | Netlify deploy preview | Automatically, on every PR |
| Production data + auth | Supabase | Migrations are applied by hand |

Merging to `main` deploys nothing. Shipping is always a deliberate act: cut a
release, and both halves go out together.

Every pull request gets a full working copy of the app — its own backend, its
own database — at no cost beyond the RAM it uses.

## How the pieces fit together

Netlify serves the site and proxies the API to the droplet:

```
https://<site>                          https://api.<domain>
├── /api/*  ──────proxied──────────────▶ backend container
└── /*      → index.html (SPA fallback)  └── Supabase (data + auth)
```

The browser only ever talks to one origin, so there is **no CORS**, and no API
hostname is baked into the JavaScript bundle. That is what lets one build run
unchanged on production and on any preview URL.

Proxy rules are generated at build time by
[frontend/scripts/netlify-redirects.sh](frontend/scripts/netlify-redirects.sh)
from `API_PROXY_TARGET`, so pointing a preview at a different backend is a
Netlify setting rather than a code change.

Locally, `frontend/vite.config.ts` proxies `/api` the same way.

The droplet itself, and how to provision one, is documented in
[deploy/README.md](deploy/README.md).

## Day-to-day

### Ship to production

Apply any pending migrations **first**. Deploys never touch Supabase — a TL
runs `make db-push` from `backend/` so schema changes get a human review, as
`backend/Makefile` intends. Release code whose schema is already in place, and
write migrations that are safe against the version still running.

Then:

```sh
git checkout main && git pull
gh release create v1.2.0 --generate-notes
```

The workflow builds the backend image from the tagged commit, pushes it to
GHCR, deploys it to the droplet, waits for `/api/v1/health` to answer, and only
then publishes the frontend. The API can always serve the new frontend before
that frontend goes live, never the reverse.

### Roll back

Actions → **Deploy production** → *Run workflow*, and give it an earlier tag.
Images are tagged by commit, so any previous build is still deployable.

Schema is not part of that. Undo a migration with a new forward migration.

The frontend can also be rolled back on its own from Netlify → **Deploys** →
pick a deploy → **Publish deploy**.

### Review a PR

Open one. Two comments appear:

- **Netlify** posts the frontend preview URL.
- **The preview workflow** posts the backend URL, `pr-<n>.<preview-domain>`.

The backend preview has its own database, migrated from that branch, isolated
from every other PR and from production. Auth comes from the non-production
Supabase project, so a preview cannot mint production sessions.

Both are destroyed when the PR closes.

To inspect a preview's data directly:

```sh
ssh deploy@<droplet> 'docker exec -it preview-db psql -U postgres -d pr_123'
```

### Point a preview frontend at its own backend

By default the Netlify preview proxies to the production API. To have it use
the PR's own backend instead, set a **deploy-preview** scoped
`API_PROXY_TARGET` in Netlify pointing at `https://pr-<n>.<preview-domain>`.

## Setup

The droplet — provisioning, DNS, Caddy, the preview database — is covered in
[deploy/README.md](deploy/README.md). What follows is the rest.

### Netlify

Connect the repository. Netlify reads [netlify.toml](netlify.toml), so build
settings configure themselves.

Set one environment variable under **Site configuration → Environment
variables**:

| Variable | Value |
| --- | --- |
| `API_PROXY_TARGET` | `https://api.<your-domain>`, no trailing slash |

Then note the **Site ID** and create a personal access token under **User
settings → Applications**.

### Repository secrets

| Secret | Where it comes from |
| --- | --- |
| `DROPLET_SSH_KEY` | The CI-only deploy key (see deploy/README.md) |
| `DROPLET_SSH_KNOWN_HOSTS` | `ssh-keyscan -t ed25519 <droplet-ip>` |
| `NETLIFY_AUTH_TOKEN` | Netlify → User settings → Applications |

### Repository variables

These are not sensitive, and the workflows read them from `vars`.

| Variable | Example |
| --- | --- |
| `DROPLET_HOST` | `164.90.x.x` |
| `DROPLET_USER` | `deploy` |
| `API_DOMAIN` | `api.inspirate.org` |
| `PREVIEW_DOMAIN` | `preview.generatenu.com` |
| `NETLIFY_SITE_ID` | From Netlify → Site configuration → General |
| `API_PROXY_TARGET` | `https://api.inspirate.org` |
| `PUBLIC_FRONTEND_URL` | `https://inspirate.org` |

### Protect production

Settings → **Environments** → `production`, with the TLs as **required
reviewers**. The deploy workflow declares `environment: production`, so every
release pauses for approval before anything ships.

## Cost

About **$12–24/month for the droplet**, shared across every Generate project
on it — not per project. Netlify's Free plan covers production hosting and
unlimited deploy previews with no commercial-use restriction. Supabase stays on
its free tier. Preview backends and preview databases cost nothing.

Two Netlify limits worth watching: 300 build minutes and 100 GB bandwidth per
month.

## Operating notes

**One droplet is a single point of failure for every project on it.** That is
downtime, not data loss — production data lives in Supabase, and the only
Postgres on the droplet holds disposable preview databases. Weekly snapshots
make a rebuild quick.

**When login is added**, note that auth reads a `jwt` cookie and the browser
talks to Netlify, not the droplet. Any `Set-Cookie` the backend sends must be
**host-only** — no `Domain` attribute pointing at the API hostname, or the
browser will reject it.

## Troubleshooting

**Deploy fails at "Wait for the backend to be healthy".** The container started
but the API is not answering. `ssh` in and check
`docker logs inspirate-consulting-backend`. The usual cause is a missing
variable in `/opt/inspirate-consulting/backend.env` — the backend exits at
startup if any `DB_*` or `SUPABASE_*` value is absent.

**Preview fails at "Bring the preview up".** Check the workflow log for which
stage failed: database creation, migrations, or the container start. A
migration referencing Supabase-managed schemas such as `auth.*` will fail here,
because the preview database is a plain Postgres.

**`docker pull` fails with `denied` on the droplet.** The GHCR login expired or
was never done. See step 5 of [deploy/README.md](deploy/README.md).

**API calls 404 in the browser but the API is up.** Check `dist/_redirects` in
the Netlify deploy log. `API_PROXY_TARGET` must have **no trailing slash**, or
the proxied path becomes `//greeting`.

**A push to `main` published the frontend.** It should not: `netlify.toml`
skips builds when `$CONTEXT` is `production`. Confirm that `ignore` line is
intact and that nobody enabled auto-publishing in the Netlify UI.

**A preview certificate fails to issue.** Confirm the wildcard DNS record for
`*.<preview-domain>` points at the droplet. Caddy issues per-hostname
certificates over HTTP-01, which needs the name to resolve before the container
starts.
