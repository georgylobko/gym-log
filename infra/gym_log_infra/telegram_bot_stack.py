"""Compute layer: a fresh Telegram webhook Lambda + public Function URL.

Everything here is created and owned by CDK (no import): the function, its IAM
execution role, and the Function URL. Deploy it independently with `cdk deploy`.

The DynamoDB table is NOT owned here — it already exists and is referenced by name
(see data_stack.py). We only grant this function access to it.

The function is self-contained (no Lambda layers). Its deps + source are bundled
into the deployment package automatically during `cdk synth/deploy` via a LOCAL
bundling hook (no Docker, no separate build step). Docker is the fallback only if
local bundling is unavailable.

Secrets (bot token, chat id) are injected as Lambda environment variables from the
repo-root .env at deploy time — never committed to git.
"""

import shutil
import subprocess
from pathlib import Path

import jsii
from aws_cdk import BundlingOptions, DockerImage, Duration, ILocalBundling, Stack
from aws_cdk import aws_lambda as lambda_
from constructs import Construct

from gym_log_infra.config import Config
from gym_log_infra.data_stack import existing_exercise_sets_table

# Bot service dir (contains requirements.txt + src/), relative to infra/.
_SERVICE_DIR = (Path(__file__).resolve().parents[1] / ".." / "services" / "telegram-bot").resolve()

# Wheels compatible with the Lambda runtime even when building on macOS.
_PLATFORM_ARGS = [
    "--python-platform", "x86_64-manylinux2014",
    "--python-version", "3.12",
    "--only-binary=:all:",
]


def _pip_install(output_dir: str) -> None:
    """Install requirements into output_dir. Prefer `uv pip` (always available in
    this project); fall back to the stdlib pip module if uv isn't on PATH."""
    if shutil.which("uv"):
        subprocess.run(
            [
                "uv", "pip", "install",
                "--target", output_dir,
                "--requirement", str(_SERVICE_DIR / "requirements.txt"),
                *_PLATFORM_ARGS,
            ],
            check=True,
        )
    else:
        subprocess.run(
            [
                "python3", "-m", "pip", "install",
                "--target", output_dir,
                "--requirement", str(_SERVICE_DIR / "requirements.txt"),
                "--platform", "manylinux2014_x86_64",
                "--python-version", "3.12",
                "--implementation", "cp",
                "--only-binary=:all:",
                "--upgrade",
            ],
            check=True,
        )


@jsii.implements(ILocalBundling)
class _LocalPipBundling:
    """Installs requirements + copies src into CDK's asset output dir, locally.

    Returning True tells CDK the asset was bundled without Docker. Any exception
    falls through to the Docker-based fallback declared in BundlingOptions.
    """

    def try_bundle(self, output_dir: str, *_args, **_kwargs) -> bool:
        _pip_install(output_dir)
        shutil.copytree(_SERVICE_DIR / "src", output_dir, dirs_exist_ok=True)
        if not (Path(output_dir) / "lambda_function.py").is_file():
            raise RuntimeError("bundling produced no lambda_function.py")
        return True


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

        # Bundle deps + source at synth/deploy time — local first, Docker fallback.
        code = lambda_.Code.from_asset(
            str(_SERVICE_DIR),
            bundling=BundlingOptions(
                image=DockerImage.from_registry("python:3.12"),
                local=_LocalPipBundling(),
                command=[
                    "bash", "-c",
                    "pip install -r requirements.txt -t /asset-output "
                    "&& cp -r src/. /asset-output",
                ],
            ),
        )

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
