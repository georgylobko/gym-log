import re
from datetime import datetime

from boto3.dynamodb.conditions import Attr
from telegram import Update
from telegram.ext import ContextTypes

from utils import to_snake_case, get_table

table = get_table()

def parse_string_to_exercise_name(message: str) -> str:
    pattern = r"^/history exercise (\w+(?: \w+)*)$"
    match = re.match(pattern, message)
    
    if match:
        return to_snake_case(match.group(1))
    else:
        return None

async def get_exercise_history_handler(update: Update, context: ContextTypes.DEFAULT_TYPE):
    exercise_name = parse_string_to_exercise_name(update.message.text)
    if not exercise_name:
        await context.bot.send_message(chat_id=update.effective_chat.id, text="Invalid format. Please use /history exercise <exercise>")
        return

    user_id = update.message.from_user.id
    response = table.scan(
        FilterExpression=Attr('name').eq(exercise_name) & Attr('user_id').eq(user_id)
    )
    items = response['Items']

    if not items:
        await context.bot.send_message(chat_id=update.effective_chat.id, text='No history found for ' + exercise_name)
        return

    items_with_date = [(item, datetime.fromisoformat(item['created_at']).date()) for item in items]
    latest_date = max(date for _, date in items_with_date)
    last_session_exercises = [item for item, date in items_with_date if date == latest_date]

    message = ''
    for item in last_session_exercises:
        message += item['created_at'].split('T')[0] + ' ' + str(item['weight']) + 'x' + str(item['reps']) + '\n'

    if not message:
        message = 'No history found for ' + exercise_name
    
    await context.bot.send_message(chat_id=update.effective_chat.id, text=message)
