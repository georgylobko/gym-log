# gym-log

A personal training-progress tracker. Log your gym sets by messaging a Telegram
bot; everything runs serverless on AWS and is managed as code with AWS CDK.

## Repository layout

```
gym-log/
├── infra/                  # AWS CDK app (Python + uv) — all AWS resources as code
├── services/
│   └── telegram-bot/       # Telegram webhook Lambda (Python)
└── frontend/               # placeholder for a future web app (S3 + CloudFront)
```

Each area has its own README:
- [`infra/README.md`](infra/README.md) — build, deploy, day-to-day ops.
- [`services/telegram-bot/README.md`](services/telegram-bot/README.md) — the bot.
- [`frontend/README.md`](frontend/README.md) — future frontend.

## How it works

```
Telegram  ──webhook POST──▶  Lambda Function URL  ──▶  telegramBotHandler (Lambda)
                                                              │
                                                              ▼
                                                    DynamoDB (exercise_sets)
```

Secrets (bot token, chat id) live in the git-ignored `.env`; CDK injects them as
Lambda environment variables at deploy time. They are never committed to the repo.

## Getting started

The bot Lambda is deployed from scratch by CDK; it reuses the existing
`exercise_sets` DynamoDB table (referenced, not managed). See
[`infra/README.md`](infra/README.md). All AWS commands use `--profile georgii`.

## ⚠️ Security

The bot token must **never** be committed. It belongs only in the git-ignored
`.env`. If a token has ever been committed or shared, rotate it in @BotFather
immediately.
