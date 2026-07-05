"""Data layer: a reference to the EXISTING DynamoDB table.

The `exercise_sets` table was created manually and holds real workout history. This
project does NOT create, own, import, or delete it — it only looks it up by name so
the bot function can be granted access. Nothing here is deployed.
"""

from aws_cdk import aws_dynamodb as dynamodb
from constructs import Construct

from gym_log_infra.config import Config


def existing_exercise_sets_table(
    scope: Construct, config: Config
) -> dynamodb.ITable:
    """Return a reference to the existing exercise_sets table (not managed by CDK)."""
    return dynamodb.Table.from_table_name(
        scope, "ExerciseSetsTable", config.data.table_name
    )
