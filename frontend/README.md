# Frontend (placeholder)

Reserved for a future gym-log web app. Not built yet.

## Planned approach

- Static SPA (framework TBD) built to `dist/`.
- Hosted on **S3 + CloudFront** (private bucket, OAC), managed by a new
  `GymLogFrontend` CDK stack in [`../infra`](../infra).
- Talks to a backend API. A Go REST backend already exists on the
  `github-actions-demo` branch (`internal/`, `sql/`) and could move to
  `services/api/` and get its own stack (Lambda/API Gateway or ECS) when promoted.

Nothing here is wired up — this file only marks where the frontend will live so the
repo structure is ready to extend.
