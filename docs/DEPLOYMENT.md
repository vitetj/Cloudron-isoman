# ISOMan Cloudron Deployment Guide

This repository is packaged as a Cloudron-native app with `proxyAuth`, `ldap` addon support, and persistence under `/app/data`.

## Prerequisites

- Cloudron instance
- Cloudron CLI installed locally

```bash
npm install -g cloudron
```

## Package Structure

- Manifest: `CloudronManifest.json`
- Container build: `Dockerfile`
- Persistent storage: Cloudron `localstorage` addon mapped to `/app/data`
- Health check endpoint: `/api/health`
- Internal app port: `8080`
- Cloudron auth wall: `proxyAuth`
- LDAP integration for app-side create auth: `ldap` addon (`CLOUDRON_LDAP_*`)

## Build and Install

From repository root:

```bash
cloudron login my.cloudron.domain
cloudron build
cloudron install --image <APP_IMAGE_ID> --location isoman.example.com
```

Or with Make targets:

```bash
make cloudron-build
make cloudron-install CLOUDRON_LOCATION=isoman.example.com
```

## Update Workflow

After code changes:

```bash
cloudron build
cloudron update
```

## If `proxyAuth` mode changed

Cloudron cannot safely toggle `proxyAuth` on an already installed app.
Use this no-data-loss workflow:

```bash
cloudron backup create --app <app-id>
cloudron uninstall --app <app-id>
cloudron install --image <APP_IMAGE_ID> --location <fqdn>
cloudron restore --app <fqdn> --backup <backup-id>
```

## Runtime Data

All runtime data is stored in:

- `/app/data/db` (SQLite database)
- `/app/data/isos` (downloaded ISO files)

## Troubleshooting

### Healthcheck failure

Check app logs:

```bash
cloudron logs --app <app-id>
```

Validate endpoint from inside container:

```bash
wget --spider -q http://127.0.0.1:8080/api/health || echo failed
```

### Login page visible but app returns 401

- Verify the app was installed with `proxyAuth` enabled (not switched dynamically)
- Check app logs: `cloudron logs --app <app-id>`
- Verify expected anonymous redirect:

```bash
curl -I https://<fqdn>/
```

Expected response is `302` to `/login`.

### Protect only `POST /api/isos`

ISOMan supports targeted protection for ISO creation with an auth challenge.

Set app environment variables:

```bash
CREATE_ISO_AUTH_ENABLED=true
```

Then either:

```bash
BASIC_AUTH_USERNAME=<user>
BASIC_AUTH_PASSWORD=<password>
```

Cloudron LDAP (recommended, `ldap` addon enabled):

```bash
LDAP_AUTH_ENABLED=true
```

External LDAP:

```bash
LDAP_AUTH_ENABLED=true
LDAP_URL=ldap://...
LDAP_USERS_BASE_DN=ou=users,dc=example,dc=com
LDAP_USER_FILTER=(|(uid={user})(mail={user}))
LDAP_BIND_DN=...
LDAP_BIND_PASSWORD=...
```
