# ISOMan (Cloudron)

ISOMan is a Cloudron-dedicated app to download, verify and serve Linux ISO files.

## Scope

This repository is intentionally minimal and focused on Cloudron packaging/runtime.

## Cloudron Deploy

```bash
npm install -g cloudron
cloudron login my.cloudron.domain
cloudron build
cloudron install --location iso.example.com
```

## Runtime

- Internal port: `8080`
- Healthcheck: `/health`
- Persistent storage: `/app/data` (Cloudron `localstorage`)

## Useful Commands

```bash
make cloudron-build
make cloudron-install CLOUDRON_LOCATION=iso.example.com
cloudron update
```

## Repository Essentials

- `CloudronManifest.json` - Cloudron app manifest
- `Dockerfile` - Cloudron runtime image
- `backend/docker-entrypoint.sh` - startup script
- `docs/DEPLOYMENT.md` - concise Cloudron operations guide

## License

MIT
