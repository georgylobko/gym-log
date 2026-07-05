import os

import boto3
from boto3.dynamodb.conditions import Attr


def to_snake_case(text: str) -> str:
    return text.lower().replace(" ", "_")


def from_snake_case(text: str) -> str:
    return text.replace("_", " ").lower()


def get_table():
    """Return the DynamoDB Table for exercise sets. Name comes from the TABLE_NAME
    env var (set by CDK), defaulting to the existing table for local runs."""
    table_name = os.environ.get("TABLE_NAME", "exercise_sets")
    return boto3.resource("dynamodb").Table(table_name)


def scan_all_user_sets(table, user_id) -> list[dict]:
    """Return ALL sets for a user, following pagination.

    A single table.scan() returns at most 1MB of data and truncates the rest via
    LastEvaluatedKey; for a full-history aggregation we must page through it all.
    """
    items: list[dict] = []
    scan_kwargs = {"FilterExpression": Attr("user_id").eq(user_id)}
    while True:
        response = table.scan(**scan_kwargs)
        items.extend(response.get("Items", []))
        last_key = response.get("LastEvaluatedKey")
        if not last_key:
            return items
        scan_kwargs["ExclusiveStartKey"] = last_key
