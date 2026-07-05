import os

import boto3


def to_snake_case(text: str) -> str:
    return text.lower().replace(" ", "_")


def from_snake_case(text: str) -> str:
    return text.replace("_", " ").lower()


def get_table():
    """Return the DynamoDB Table for exercise sets. Name comes from the TABLE_NAME
    env var (set by CDK), defaulting to the existing table for local runs."""
    table_name = os.environ.get("TABLE_NAME", "exercise_sets")
    return boto3.resource("dynamodb").Table(table_name)
