#!/usr/bin/env python3
"""CDK app entrypoint for gym-log.

One deployable stack:
  - GymLogTelegramBot : a fresh Lambda + public Function URL for the Telegram
    webhook, granted access to the EXISTING (unmanaged) exercise_sets table.

The DynamoDB table is referenced by name, never created or deleted here.
"""

import aws_cdk as cdk

from gym_log_infra.config import CONFIG
from gym_log_infra.telegram_bot_stack import GymLogTelegramBotStack

app = cdk.App()
env = cdk.Environment(account=CONFIG.account, region=CONFIG.region)

GymLogTelegramBotStack(app, "GymLogTelegramBot", config=CONFIG, env=env)

app.synth()
