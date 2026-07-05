"""/summary [months] — aggregate the user's training data over the last N months
(default 6) into an LLM-ready JSON summary, intended as input for Bedrock analysis.
"""

import io
import json
import re
from datetime import datetime, timezone

from telegram import Update
from telegram.ext import ContextTypes

from utils import get_table, scan_all_user_sets
from aggregation import aggregate_training_data

table = get_table()

# Telegram hard-caps text messages at 4096 chars; send larger payloads as a file.
_MAX_TELEGRAM_TEXT = 3500


def parse_months(message: str) -> int | None:
    """Parse '/summary' or '/summary <months>'. Returns the month count (default 6),
    or None if the message doesn't match."""
    match = re.match(r"^/summary(?:\s+(\d{1,2}))?$", message.strip())
    if not match:
        return None
    return int(match.group(1)) if match.group(1) else 6


async def summary_handler(update: Update, context: ContextTypes.DEFAULT_TYPE):
    months = parse_months(update.message.text)
    if months is None or months < 1:
        await context.bot.send_message(
            chat_id=update.effective_chat.id,
            text="Invalid format. Use /summary or /summary <months> (e.g. /summary 6)",
        )
        return

    user_id = update.message.from_user.id
    items = scan_all_user_sets(table, user_id)

    # created_at is written as a naive UTC timestamp; compare against naive UTC now.
    now = datetime.now(timezone.utc).replace(tzinfo=None)
    result = aggregate_training_data(items, now=now, months=months)

    if result["summary"]["total_sets"] == 0:
        await context.bot.send_message(
            chat_id=update.effective_chat.id,
            text=f"No training data found in the last {months} months.",
        )
        return

    payload = json.dumps(result, indent=2, default=str)

    if len(payload) <= _MAX_TELEGRAM_TEXT:
        await context.bot.send_message(
            chat_id=update.effective_chat.id, text=f"```json\n{payload}\n```", parse_mode="Markdown"
        )
    else:
        # Too big for one message — send as a .json document.
        doc = io.BytesIO(payload.encode("utf-8"))
        doc.name = f"training_summary_{months}m.json"
        await context.bot.send_document(
            chat_id=update.effective_chat.id,
            document=doc,
            caption=f"Training summary — last {months} months "
            f"({result['summary']['total_sets']} sets, "
            f"{result['summary']['distinct_exercises']} exercises).",
        )
