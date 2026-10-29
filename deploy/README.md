# inspirate-consulting on the shared droplet

This project's backend and its per-PR previews run on the DigitalOcean droplet
that several Generate projects share. The droplet and its shared stack are
documented in [GenerateNU/droplet](https://github.com/GenerateNU/droplet); this
directory is only what inspirate-consulting adds on top.

**None of this project's irreplaceable data lives here.** Production data and
auth stay in Supabase, and its preview databases are disposable. If the droplet
is lost, it is rebuilt from GenerateNU/droplet and this directory. You lose
uptime, not data.

```
droplet
├── caddy          TLS + routes by hostname, driven by container labels
├── preview-db     Postgres: preview databases only, disposable
├── <project>-backend          →  api.<project>.org
└── <project>-pr-<n>           →  pr-<n>.<preview-domain>
```

Caddy is `caddy-docker-proxy`: it reads the Docker socket and builds its config
from container labels. A preview container starts with a `caddy` label and
routes itself, certificate included. Nothing edits a Caddyfile, and nothing
reloads.

## Provisioning

The droplet itself (Docker, the deploy user, firewall, GHCR login, and the
shared Caddy + preview-db stack) is set up once for every project, and is
documented in [GenerateNU/droplet](https://github.com/GenerateNU/droplet).
What follows is only this project's part.

### 1. DNS

| Record | Points at | Purpose |
| --- | --- | --- |
| `api.<project>.org` | droplet IP | Production backend |
| `*.<preview-domain>` | droplet IP | Every PR preview hostname |

### 2. Set the project up

Once per project:

```sh
sudo install -o deploy -g deploy -d /opt/inspirate-consulting
cd /opt/inspirate-consulting
# from deploy/app/, fill in and copy:
#   backend.env.example  -> backend.env
#   preview.env.example  -> preview.env
chmod 600 backend.env preview.env

# The preview domain, read by preview-up.sh:
echo "preview.generatenu.com" > preview-domain
```

`DB_PASSWORD` in `preview.env` must match `PREVIEW_DB_PASSWORD` in
`/opt/shared/.env` (see GenerateNU/droplet) — the preview containers connect to that Postgres as
superuser, and `preview-up.sh` reads the password from this file.

The deploy workflow copies `docker-compose.yml` here on every run, so it does
not need to be placed by hand.

### 3. Repository secrets and variables

Secrets: `DROPLET_SSH_KEY` and `DROPLET_SSH_KNOWN_HOSTS`.

Variables: `DROPLET_HOST`, `DROPLET_USER` (`deploy`), `API_DOMAIN`,
`PREVIEW_DOMAIN`.

Until `DROPLET_HOST`, `API_DOMAIN` and `PREVIEW_DOMAIN` are set, deploys and
previews cannot reach the droplet.

## Day-to-day

**Deploy the backend.** Publishing a GitHub Release does it. Apply any pending
migrations first — `make db-push` from `backend/`, reviewed by a TL; deploys do
not touch Supabase.

The workflow builds the image from the tagged commit, pushes it to GHCR, copies
the compose file up, restarts the container, and polls `/api/v1/health` before
letting the frontend go live. See [DEPLOYMENT.md](../DEPLOYMENT.md).

**Roll back.** Actions → *Deploy production* → Run workflow with an earlier
tag. Images are tagged by commit, so any past build is still deployable. Schema
is not part of that — undo a migration with a new forward migration.

**Previews** happen on every PR with no action required. The bot comments the
API URL; the container and database are removed when the PR closes.

**Inspect a preview database.**

```sh
ssh deploy@<droplet-ip>
docker exec -it preview-db psql -U postgres -d pr_123
```

**Clean up an orphan** (cleanup failed, or a PR predates this setup):

```sh
ssh deploy@<droplet-ip> bash -s -- 123 < deploy/preview-down.sh
```

## Operating notes

**This is a single point of failure for every project on it.** A full disk or a
bad `docker compose up` takes all of them down at once. That is downtime, not
data loss — but it is still four clients at once. Enable weekly snapshots
(+20% of droplet cost) so a rebuild is quick.

**Watch disk.** Every project on the droplet shares one disk. On 2026-10-26
it filled with old images and every deploy failed. Deploys now remove this
project's superseded images, and the droplet prunes unused ones nightly (see
GenerateNU/droplet). Never run `docker system prune --volumes` there: it
deletes every project's preview databases, and other projects' production
data.

**Keep OS packages current.** `unattended-upgrades` is the cheap answer:

```sh
apt install unattended-upgrades && dpkg-reconfigure -plow unattended-upgrades
```

Adding another project to the droplet is covered in
[GenerateNU/droplet](https://github.com/GenerateNU/droplet).
