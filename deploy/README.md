# Shared droplet

One DigitalOcean droplet hosting the backends for several Generate projects,
plus a per-PR preview backend for each of them.

**Nothing irreplaceable lives here.** Production data and auth stay in
Supabase; the only Postgres on the droplet holds preview databases, which are
disposable by definition. If the droplet is lost you rebuild it from this
directory — you lose uptime, not data. That is the whole reason this is a
reasonable thing for a student org to own.

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

### 1. Create the droplet

A **2 GB / 1 vCPU** basic droplet ($12/month) fits four Go backends, the
preview Postgres, Caddy, and a handful of preview containers. Go up to 4 GB
($24/month) if you expect many concurrent previews.

Pick Ubuntu LTS and add your team's SSH keys.

### 2. Install Docker

```sh
ssh root@<droplet-ip>
curl -fsSL https://get.docker.com | sh
```

Create a non-root deploy user so CI is not connecting as root:

```sh
adduser --disabled-password --gecos "" deploy
usermod -aG docker deploy
install -o deploy -g deploy -m 700 -d /home/deploy/.ssh
```

### 3. Give CI a deploy key

On your machine, generate a key **used only by CI**:

```sh
ssh-keygen -t ed25519 -f droplet-deploy-key -C "github-actions" -N ""
```

Append the public half to `/home/deploy/.ssh/authorized_keys` on the droplet,
then capture the host key for pinning:

```sh
ssh-keyscan -t ed25519 <droplet-ip>
```

Store the private key as the `DROPLET_SSH_KEY` secret and the `ssh-keyscan`
output as `DROPLET_SSH_KNOWN_HOSTS`. Pinning matters: without it CI would
accept any host presenting itself at that address, deploy key included.

### 4. Firewall

```sh
ufw allow OpenSSH
ufw allow 80
ufw allow 443
ufw enable
```

Postgres is deliberately **not** in that list. It publishes no ports and is
reachable only from other containers on the `edge` network.

### 5. Let the droplet pull from GHCR

Images are pushed to GitHub Container Registry, and GHCR packages are private
by default even when the repository is public. The droplet must authenticate
or the first deploy fails with `denied`.

Create a classic PAT with **only** `read:packages`, then on the droplet:

```sh
echo "<token>" | docker login ghcr.io -u <github-username> --password-stdin
```

That writes `~/.docker/config.json` for the `deploy` user and persists across
reboots. Use a token that can do nothing but read packages — it lives in
plaintext on the droplet.

Alternatively, make the two packages public (repo → Packages → Package
settings → Change visibility) and skip the login entirely. Fine for
open-source projects; not for client work where the image may embed config.

### 6. DNS

| Record | Points at | Purpose |
| --- | --- | --- |
| `api.<project>.org` | droplet IP | Production backend |
| `*.<preview-domain>` | droplet IP | Every PR preview hostname |

The wildcard is what lets a new preview hostname work the instant its
container starts. Caddy then issues a certificate per hostname over HTTP-01,
which needs no DNS credentials.

> Let's Encrypt allows 50 certificates per registered domain per week. Heavy PR
> churn across four projects could approach that; if it becomes a problem,
> switch Caddy to a DNS-challenge wildcard certificate.

### 7. Bring up the shared stack

Once per droplet:

```sh
ssh deploy@<droplet-ip>
sudo install -o deploy -g deploy -d /opt/shared
# copy deploy/shared/docker-compose.yml to /opt/shared/
cp .env.example .env   # fill in ACME_EMAIL and PREVIEW_DB_PASSWORD
chmod 600 .env
docker compose up -d
```

`docker network ls` should now show `edge`.

### 8. Set the project up

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
`/opt/shared/.env` — the preview containers connect to that Postgres as
superuser, and `preview-up.sh` reads the password from this file.

The deploy workflow copies `docker-compose.yml` here on every run, so it does
not need to be placed by hand.

### 9. Repository secrets and variables

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

**Watch disk.** Preview images accumulate in the local Docker cache. The deploy
workflow runs `docker image prune -f`, but a monthly
`docker system prune -af --volumes` is worth doing — note that `--volumes`
would also drop preview databases, so only run it when no previews are open.

**Keep OS packages current.** `unattended-upgrades` is the cheap answer:

```sh
apt install unattended-upgrades && dpkg-reconfigure -plow unattended-upgrades
```

## Adding another project

1. Copy `deploy/app/` into the other repo and change `name:` and
   `container_name:` in `docker-compose.yml`.
2. Copy the three workflows and change `PROJECT` in `preview-up.sh` and
   `preview-down.sh`.
3. `/opt/<project>/` on the droplet with its own `backend.env`, `preview.env`,
   and `preview-domain`.
4. DNS for its API hostname and preview wildcard.
5. Its own repository secrets and variables.

The shared stack is untouched — that is the point.
