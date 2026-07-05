import re

from boto3.dynamodb.conditions import Attr
from telegram import Update
from telegram.ext import ContextTypes

from utils import to_snake_case, get_table

table = get_table()

def parse_string_to_exercise_name(message: str) -> str:
    pattern = r"^/search (\w+(?: \w+)*)$"
    match = re.match(pattern, message)
    
    if match:
        return to_snake_case(match.group(1))
    else:
        return None

async def search_exercise_handler(update: Update, context: ContextTypes.DEFAULT_TYPE):
    exercise_name = parse_string_to_exercise_name(update.message.text)
    if not exercise_name:
        await context.bot.send_message(chat_id=update.effective_chat.id, text="Invalid format. Please use /search <exercise>")
        return

    user_id = update.message.from_user.id
    response = table.scan(
        FilterExpression=Attr('name').contains(exercise_name) & Attr('user_id').eq(user_id)
    )
    names = {item['name'] for item in response['Items']}

    if not names:
        await context.bot.send_message(chat_id=update.effective_chat.id, text='No exercise found')
        return

    for name in names:
        await context.bot.send_message(chat_id=update.effective_chat.id, text=str(name))
