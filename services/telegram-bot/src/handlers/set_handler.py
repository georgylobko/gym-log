import re
import uuid
from datetime import datetime
from typing import Optional, TypedDict

from telegram import Update
from telegram.ext import ContextTypes

from utils import to_snake_case, get_table

table = get_table()

class ExerciseSet(TypedDict):
    exercise: str
    weight: int
    reps: int

def parse_string_to_exercise_set(message: str) -> Optional[ExerciseSet]:
    pattern = r"^/set (\w+(?: \w+)*) (\d+) (\d+)$"
    match = re.match(pattern, message)
    
    if match:
        exercise = match.group(1)
        weight = int(match.group(2))
        reps = int(match.group(3))
        return ExerciseSet(
            exercise=exercise,
            weight=weight,
            reps=reps
        )
    else:
        return None

async def set_handler(update: Update, context: ContextTypes.DEFAULT_TYPE):
    exercise_set = parse_string_to_exercise_set(update.message.text)
    if not exercise_set:
        await context.bot.send_message(chat_id=update.effective_chat.id, text="Invalid format. Please use /set <exercise> <weight> <reps>")
        return

    item = {
        'id': str(uuid.uuid4()),
        'reps': exercise_set["reps"],
        'weight': exercise_set["weight"],
        'name': to_snake_case(exercise_set["exercise"]),
        'created_at': datetime.utcnow().isoformat(),
        'user_id': update.message.from_user.id
    }
    table.put_item(Item=item)


    await context.bot.set_message_reaction(chat_id=update.effective_chat.id, message_id=update.message.message_id, reaction="🔥")
