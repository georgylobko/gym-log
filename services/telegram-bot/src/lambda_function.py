import json
import os
import asyncio

from telegram import Update
from telegram.ext import ApplicationBuilder, CommandHandler, MessageHandler, filters

from handlers.set_handler import set_handler
from handlers.start_handler import start_handler
from handlers.get_exercise_history_handler import get_exercise_history_handler
from handlers.search_exercise_handler import search_exercise_handler
from handlers.summary_handler import summary_handler


def get_bot_token() -> str:
    """Read the bot token from the TOKEN env var (set by CDK from .env)."""
    token = os.environ.get("TOKEN")
    if not token:
        raise RuntimeError("TOKEN env var is not set.")
    return token


# Build the application and register handlers ONCE at module load. Registering
# inside the request handler (as before) duplicated handlers on warm containers,
# making each command fire multiple times.
application = ApplicationBuilder().token(get_bot_token()).build()
application.add_handler(CommandHandler("start", start_handler))
application.add_handler(
    MessageHandler(filters.TEXT & filters.Regex(r"^/set (\w+(?: \w+)*) (\d+) (\d+)$"), set_handler)
)
application.add_handler(
    MessageHandler(
        filters.TEXT & filters.Regex(r"^/history exercise (\w+(?: \w+)*)$"),
        get_exercise_history_handler,
    )
)
application.add_handler(
    MessageHandler(filters.TEXT & filters.Regex(r"^/search (\w+(?: \w+)*)$"), search_exercise_handler)
)
application.add_handler(CommandHandler("summary", summary_handler))


def lambda_handler(event, context):
    return asyncio.get_event_loop().run_until_complete(main(event, context))


async def main(event, context):
    try:
        await application.initialize()
        await application.process_update(
            Update.de_json(json.loads(event["body"]), application.bot)
        )
        return {"statusCode": 200, "body": "Success"}
    except Exception:
        return {"statusCode": 500, "body": "Failure"}
