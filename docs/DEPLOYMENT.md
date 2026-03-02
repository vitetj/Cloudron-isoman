# ISOMan Cloudron Deployment Guide

This repository is packaged for Cloudron deployment.

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
- Health check endpoint: `/health`
- Internal app port: `8080`

## Build and Install

From repository root:

```bash
cloudron login my.cloudron.domain
cloudron build
cloudron install --location isoman.example.com
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

## Runtime Data

All runtime data is stored in:

- `/app/data/db` (SQLite database)
- `/app/data/isos` (downloaded ISO files)

## Troubleshooting

### Entrypoint error (`/entrypoint.sh: no such file or directory`)

The image build normalizes shell script line endings to LF. Rebuild with:

```bash
cloudron build --no-cache
```

### Healthcheck failure

Check app logs:

```bash
cloudron logs --app <app-id>
```

Validate endpoint from inside container:

```bash
wget --spider -q http://127.0.0.1:8080/health || echo failed
```
