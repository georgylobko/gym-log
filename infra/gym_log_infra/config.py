"""Central configuration for the gym-log infrastructure.

The Lambda + Function URL + role are created fresh by CDK (deployed independently).
The DynamoDB table already exists and is only *referenced* (by name) — never
created or deleted by this app.

AWS account/region come from the repo-root .env file (see .env.example). The .env
values are authoritative here: they override any AWS_* vars already in the shell,
so `cdk synth/deploy` always targets the project's account regardless of ambient
environment.
"""

import os
from dataclasses import dataclass, field
from pathlib import Path


def _load_dotenv() -> None:
    """Load the repo-root .env (KEY=VALUE lines) into os.environ, OVERRIDING any
    existing values so .env is the single source of truth. Zero-dependency."""
    # config.py -> gym_log_infra -> infra -> repo root
    env_path = Path(__file__).resolve().parents[2] / ".env"
    if not env_path.is_file():
        return
    for raw in env_path.read_text().splitlines():
        line = raw.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        key, _, value = line.partition("=")
        os.environ[key.strip()] = value.strip().strip("'\"")


def _require(name: str) -> str:
    value = os.environ.get(name)
    if not value:
        raise RuntimeError(
            f"{name} is not set. Add it to the repo-root .env file "
            f"(see .env.example) or export it in your shell."
        )
    return value


_load_dotenv()

# Deployment target — sourced from .env (AWS_ACCOUNT_ID / AWS_REGION).
ACCOUNT = _require("AWS_ACCOUNT_ID")
REGION = _require("AWS_REGION")

# Telegram secrets — sourced from .env (TOKEN / CHATID) and set as Lambda env vars
# at deploy time. NOTE: env vars are stored in the (git-ignored) CloudFormation
# template and visible in the Lambda console — weaker than SSM, but matches the
# original setup and keeps secrets out of git.
TELEGRAM_TOKEN = _require("TOKEN")
TELEGRAM_CHAT_ID = _require("CHATID")


@dataclass(frozen=True)
class DataConfig:
    # Physical name of the EXISTING DynamoDB table (referenced, not managed).
    table_name: str = "exercise_sets"


@dataclass(frozen=True)
class BotConfig:
    # Name for the NEW, CDK-created function. Distinct from the old manual function
    # ("telegramBotHandler") so the two never collide.
    function_name: str = "gym-log-telegram-bot"
    handler: str = "lambda_function.lambda_handler"
    memory_mb: int = 256
    timeout_seconds: int = 10
    # Dependencies (python-telegram-bot, requests) are bundled into the deployment
    # package from services/telegram-bot/requirements.txt — no Lambda layers.


@dataclass(frozen=True)
class Config:
    account: str = ACCOUNT
    region: str = REGION
    env_name: str = "prod"
    data: DataConfig = field(default_factory=DataConfig)
    bot: BotConfig = field(default_factory=BotConfig)
    telegram_token: str = TELEGRAM_TOKEN
    telegram_chat_id: str = TELEGRAM_CHAT_ID


CONFIG = Config()
