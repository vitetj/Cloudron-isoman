# ISOMan - Cloudron Native Package

This repository packages ISOMan as a Cloudron-native app with Cloudron `proxyAuth`, persistent storage in `/app/data`, and runtime on `cloudron/base:5`.

## Package Layout

- `CloudronManifest.json`: Cloudron manifest (port, addons, healthcheck)
- `CloudronVersions.json`: Cloudron versions index for publishing
- `Dockerfile`: multi-stage build (Bun UI + Go backend) + Cloudron runtime
- `start.sh`: Cloudron runtime bootstrap (`PORT`, `DATA_DIR`, persistent dirs)
- `backend/`: Go backend source + migrations
- `ui/`: React/RSBuild frontend source

## Technical Decisions

- Fixed internal port: `8080`
- Persistent data path: `/app/data`
- Healthcheck endpoint: `/api/health`
- Global authentication: Cloudron `proxyAuth`
- Optional create-only auth (`POST /api/isos`): Cloudron LDAP (`ldap` addon) or Basic fallback
- Cloudron app icon: `appicon.png` (PNG for login screen compatibility)

## Build & Install (Cloudron)

```bash
npm install -g cloudron
cloudron login my.cloudron.domain
cloudron build
cloudron install --image <APP_IMAGE_ID> --location iso.example.com
```

To update an existing install:

```bash
cloudron update
```

## Local Docker Validation

```bash
docker build -t isoman-cloudron .
docker run --rm -p 8080:8080 -v %cd%/cloudron-data:/app/data isoman-cloudron
```

Validate:
- UI at `http://localhost:8080`
- Health API at `http://localhost:8080/api/health`
- WebSocket endpoint at `/ws`
- DB + ISO persistence under `./cloudron-data`

## Publishing

Initialize versions file (if needed):

```bash
cloudron versions init
```

Add a built version:

```bash
cloudron versions add --version $(cat VERSION) --docker-image <REGISTRY/IMAGE:TAG>
```

Publish using your regular Cloudron workflow (git repo, artifacts, submission).

## Cloudron Auth Notes

The manifest enables both `proxyAuth` (global access wall) and `ldap` (injects `CLOUDRON_LDAP_*` vars).

Important: `proxyAuth` cannot be added dynamically to an already-installed app. If the app was installed without it, use backup → uninstall → reinstall → restore.

## Troubleshooting

### Login page shown but app returns `401`

- Ensure the app was installed with `proxyAuth` enabled (not toggled later)
- Check logs: `cloudron logs --app <app-id>`
- Verify anonymous redirect:

```bash
curl -I https://<fqdn>/
```

Expected response: `302` to `/login`.

### ISO create returns HTML/non-JSON errors in browser

The frontend now handles auth challenges on `POST /api/isos` and prompts for credentials when needed. If needed, hard-refresh the UI (`Ctrl+F5`) after updates.

## Optional Create-Only Protection

Protect only ISO creation (`POST /api/isos`) with authentication challenge.

Cloudron env:

- `CREATE_ISO_AUTH_ENABLED=true`
- Option 1 (simple): `BASIC_AUTH_USERNAME` + `BASIC_AUTH_PASSWORD`
- Option 2 (recommended on Cloudron): `LDAP_AUTH_ENABLED=true` (uses injected `CLOUDRON_LDAP_*`)
- Option 3 (external LDAP): `LDAP_URL`, `LDAP_USERS_BASE_DN`, `LDAP_USER_FILTER`, and optional `LDAP_BIND_DN` / `LDAP_BIND_PASSWORD`

Behavior:

- `POST /api/isos` without valid credentials returns `401`
- Other routes (`GET /api/isos`, UI, etc.) remain under standard Cloudron `proxyAuth`
