# Brick Hill Website V1 — Render Free

This version is prepared for a Render Free **Web Service**.

## Render settings

If you create the service manually:

- **Runtime:** Node
- **Build Command:** `npm install`
- **Start Command:** `npm start`
- **Root Directory:** leave empty
- **Health Check Path:** `/health`
- **Instance Type:** Free

The included `render.yaml` contains the same configuration.

## Deploy

1. Upload/push this folder to GitHub.
2. In Render choose **New → Web Service**.
3. Connect the repository.
4. Use the settings above, or let the Blueprint use `render.yaml`.
5. Deploy.

The server listens on `0.0.0.0` and uses Render's `PORT` environment variable, so it is suitable for Render Web Services.

## Important Free-plan limitation

The current V1 stores demo users and games in `data/site.json`. Render Free web services have an **ephemeral filesystem**, so changes to local files can disappear after a restart, redeploy, or spin-down. This means V1 is suitable for testing, but its account/game data is not production-persistent.

For a real public Brick Hill site, the next step is to move users, games, inventory and other persistent data to a database. Render documents Postgres as the persistent database option, while its Free Postgres offering currently has a 30-day expiration.

Free web services also spin down after 15 minutes without inbound traffic, so the first request after sleeping can take around a minute.

## Local run

    npm install
    npm start

Then open:

    http://localhost:3000
