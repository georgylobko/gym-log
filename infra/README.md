# gym-log infrastructure (CDK)

AWS CDK app (Python, managed with [uv](https://docs.astral.sh/uv/)) that deploys
the gym-log Telegram bot.

One deployable stack:

| Stack | Resources | Ownership |
| --- | --- | --- |
| `GymLogTelegramBot` | Lambda `gym-log-telegram-bot`, its IAM role, a public Function URL | **created & owned by CDK** |
| — | DynamoDB `exercise_sets` table | **pre-existing**, only *referenced* (never created/deleted) |

The function is deployed **from scratch** and independently. Its dependencies +
source are built into `services/telegram-bot/dist/` and shipped by CDK — no Lambda
layers, no Docker. The existing `exercise_sets` table is looked up by name and the
function is granted access to it.

Account/region/profile come from the repo-root `.env` (see `.env.example`). All
commands below use `--profile georgii`. `aws-cdk-lib` is pinned in `pyproject.toml`
to match a locally-installed `cdk` CLI (v2.178.x), so the `cdk` command runs
directly — no Node/npm needed.

---

## One-time setup

```bash
cd infra
uv sync        # Python deps (aws-cdk-lib, pinned to match your cdk CLI)
```
Requires the `cdk` CLI installed (v2.178.x). If yours differs, see the version note
above.

Set secrets in the repo-root `.env` (copy from `.env.example`). CDK reads `TOKEN`
and `CHATID` at deploy time and sets them as the Lambda's env vars. `.env` is
git-ignored — never commit it. Rotate the token in @BotFather if it was ever shared.

Bootstrap CDK (once per account/region):
```bash
cdk bootstrap --profile georgii
```

---

## Deploy

```bash
# 1. Build the Lambda package (deps + source -> services/telegram-bot/dist)
../services/telegram-bot/build.sh

# 2. Review, then deploy the fresh stack
cdk diff   --profile georgii
cdk deploy --profile georgii
```

`cdk deploy` prints the **Function URL**. Point Telegram's webhook at it (one time,
and again whenever the URL changes):
```bash
curl "https://api.telegram.org/bot<BOT_TOKEN>/setWebhook?url=<FUNCTION_URL>"
```

Verify in Telegram: `/start`, then `/set bench 60 5`.

---

## Everyday workflow

Edit code in [`../services/telegram-bot/src`](../services/telegram-bot), then:
```bash
cd infra
uv run pytest                        # synth-level guard tests
../services/telegram-bot/build.sh    # rebuild the package
cdk deploy --profile georgii
```

---

## Notes & safety

- **The table is never touched by CDK.** It's referenced by name; `cdk destroy`
  removes only the function/role/URL, never your workout history.
- Always `cdk diff` before `deploy`.
- Secrets live only in the git-ignored `.env`. They ARE set as Lambda env vars (and
  thus appear in the git-ignored `cdk.out/` template and the Lambda console) — this
  is weaker than a secrets manager but keeps them out of git. Never commit `.env`.
- The old manual function (`telegramBotHandler`) is separate. Once this new bot
  works, delete the old function + its webhook to avoid double-processing updates.

## Future

- **GitHub Actions CI/CD** on push to `main` (deferred; GitHub OIDC role, no keys).
- **Frontend**: S3 + CloudFront static site as a new stack — see
  [`../frontend/README.md`](../frontend/README.md).
