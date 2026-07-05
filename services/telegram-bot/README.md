# Telegram bot service

AWS Lambda that receives Telegram webhook updates and logs gym exercise sets to
DynamoDB. Invoked via a public Lambda Function URL that Telegram posts to.

## Commands

| Command | Action |
| --- | --- |
| `/start` | Health-check reply |
| `/set <exercise> <weight> <reps>` | Store a set, react 🔥 |
| `/history exercise <exercise>` | Show your most recent session for that exercise |
| `/search <exercise>` | List exercise names matching the query |
| `/summary [months]` | Aggregate training data over the last N months (default 6) into an LLM-ready JSON summary — intended as Bedrock input |

## Layout

- `src/lambda_function.py` — entrypoint; builds the PTB application and registers
  handlers once at module load. Reads the bot token from the `TOKEN` env var.
- `src/handlers/` — one module per command.
- `src/aggregation.py` — pure (Telegram/boto3-free) training-data aggregation used
  by `/summary`; reusable directly by a future Bedrock analysis pipeline.
- `src/utils.py` — helpers; `get_table()` reads `TABLE_NAME`, `scan_all_user_sets()`
  pages through all of a user's sets.
- `requirements.txt` — runtime deps, bundled into the package (no Lambda layers).
- `build.sh` — installs deps + copies `src/` into `dist/`; CDK ships `dist/`.
- `dist/` — build output (git-ignored). Regenerate with `./build.sh`.

## Configuration (env vars set by CDK)

| Var | Purpose |
| --- | --- |
| `TABLE_NAME` | DynamoDB table name |
| `TOKEN` | Telegram bot token |
| `CHATID` | Allowed Telegram chat id |

## Deploy

Build, then deploy via CDK:
```bash
./build.sh                                    # deps + src -> dist/
cd ../../infra && cdk deploy --profile georgii
```
Full instructions: [`../../infra/README.md`](../../infra/README.md).
