"""Compute layer: a fresh Telegram webhook Lambda + public Function URL.

Everything here is created and owned by CDK (no import): the function, its IAM
execution role, and the Function URL. Deploy it independently with `cdk deploy`.

The DynamoDB table is NOT owned here — it already exists and is referenced by name
(see data_stack.py). We only grant this function access to it.

The function is self-contained (no Lambda layers): its deps + source are built
into services/telegram-bot/dist by build.sh, and CDK ships that directory. Run
the build before `cdk synth/deploy` (see infra/README.md).

Secrets (bot token, chat id) are injected as Lambda environment variables from the
repo-root .env at deploy time — never committed to git.
"""

from pathlib import Path

from aws_cdk import Duration, Stack
from aws_cdk import aws_lambda as lambda_
from constructs import Construct

from gym_log_infra.config import Config
from gym_log_infra.data_stack import existing_exercise_sets_table

# Pre-built deployment directory (deps + source), relative to infra/.
_DIST = Path(__file__).resolve().parents[1] / ".." / "services" / "telegram-bot" / "dist"


class GymLogTelegramBotStack(Stack):
    def __init__(
        self,
        scope: Construct,
        construct_id: str,
        config: Config,
        **kwargs,
    ) -> None:
        super().__init__(scope, construct_id, **kwargs)

        bot = config.bot
        exercise_sets_table = existing_exercise_sets_table(self, config)

        # Ship the pre-built deployment directory (deps + source). No Docker.
        if not (_DIST / "lambda_function.py").is_file():
            raise FileNotFoundError(
                f"{_DIST} is not built. Run services/telegram-bot/build.sh before "
                f"cdk synth/deploy (see infra/README.md)."
            )
        code = lambda_.Code.from_asset(str(_DIST))

        # Fresh function with a CDK-managed execution role (created automatically).
        self.function = lambda_.Function(
            self,
            "TelegramBotFunction",
            function_name=bot.function_name,
            runtime=lambda_.Runtime.PYTHON_3_12,
            handler=bot.handler,
            code=code,
            memory_size=bot.memory_mb,
            timeout=Duration.seconds(bot.timeout_seconds),
            environment={
                "TABLE_NAME": config.data.table_name,
                "TOKEN": config.telegram_token,
                "CHATID": config.telegram_chat_id,
            },
        )

        # Public Function URL that Telegram posts webhook updates to. This is a NEW
        # URL — re-point the webhook after the first deploy (see infra/README.md).
        self.function_url = self.function.add_function_url(
            auth_type=lambda_.FunctionUrlAuthType.NONE,
            invoke_mode=lambda_.InvokeMode.BUFFERED,
        )

        # Access to the existing table.
        exercise_sets_table.grant_read_write_data(self.function)
