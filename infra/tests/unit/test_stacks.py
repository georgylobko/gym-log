"""Synthesis tests for the fresh Telegram bot stack.

Guards:
  - the function is created with the expected config,
  - a public Function URL exists,
  - the DynamoDB table is only *referenced* (never created/deleted here),
  - the Telegram secrets are passed as Lambda env vars.
"""

import aws_cdk as cdk
from aws_cdk.assertions import Match, Template

from gym_log_infra.config import CONFIG
from gym_log_infra.telegram_bot_stack import GymLogTelegramBotStack


def _synth() -> Template:
    app = cdk.App()
    env = cdk.Environment(account=CONFIG.account, region=CONFIG.region)
    stack = GymLogTelegramBotStack(app, "GymLogTelegramBot", config=CONFIG, env=env)
    return Template.from_stack(stack)


def test_lambda_has_expected_config():
    _synth().has_resource_properties(
        "AWS::Lambda::Function",
        {
            "FunctionName": CONFIG.bot.function_name,
            "Handler": CONFIG.bot.handler,
            "Runtime": "python3.12",
            "MemorySize": CONFIG.bot.memory_mb,
            "Timeout": CONFIG.bot.timeout_seconds,
        },
    )


def test_public_function_url_exists():
    _synth().has_resource_properties(
        "AWS::Lambda::Url", {"AuthType": "NONE", "InvokeMode": "BUFFERED"}
    )


def test_execution_role_is_created():
    # A fresh (CDK-owned) execution role, not an imported ARN.
    _synth().resource_count_is("AWS::IAM::Role", 1)


def test_table_is_referenced_not_created():
    # The existing table must never be created or deleted by this stack.
    _synth().resource_count_is("AWS::DynamoDB::Table", 0)


def test_lambda_env_has_secrets_and_table():
    _synth().has_resource_properties(
        "AWS::Lambda::Function",
        {
            "Environment": {
                "Variables": Match.object_equals(
                    {
                        "TABLE_NAME": CONFIG.data.table_name,
                        "TOKEN": CONFIG.telegram_token,
                        "CHATID": CONFIG.telegram_chat_id,
                    }
                )
            }
        },
    )
