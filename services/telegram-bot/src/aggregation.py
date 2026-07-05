"""Pure aggregation of exercise-set records into an LLM-friendly training summary.

No Telegram, no boto3 — just data in, structured dict out. This is intentionally
decoupled so a future Bedrock analysis pipeline can import and reuse it directly.

Input: an iterable of set records, each a dict with:
    name        str    exercise name (snake_case)
    weight      number kg (or lb) lifted
    reps        number repetitions
    created_at  str    naive ISO-8601 timestamp (as written by set_handler)

Output (see aggregate_training_data): a JSON-serializable dict summarizing the
period, an overall summary, a monthly time series, and per-exercise progression.
"""

from collections import defaultdict
from datetime import datetime, timedelta
from typing import Iterable, Optional


# Average days per month — the lookback window is approximate, which is fine for
# a "last N months" training summary.
_DAYS_PER_MONTH = 30.44


def estimated_1rm(weight: float, reps: int) -> float:
    """Epley one-rep-max estimate: weight * (1 + reps/30). A single heavy-ish set
    is a reasonable strength proxy for trend analysis."""
    return weight * (1 + reps / 30.0)


def _round(value: float, digits: int = 1) -> float:
    return round(float(value), digits)


def _session_stats(sets: list[dict]) -> dict:
    """Summarize the sets performed for one exercise on one day."""
    volume = sum(s["weight"] * s["reps"] for s in sets)
    top_1rm = max(estimated_1rm(s["weight"], s["reps"]) for s in sets)
    return {
        "sets": len(sets),
        "total_reps": sum(s["reps"] for s in sets),
        "volume": _round(volume),
        "max_weight": max(s["weight"] for s in sets),
        "est_1rm": _round(top_1rm),
    }


def _pct_change(first: float, last: float) -> Optional[float]:
    """Percent change from first to last; None if no meaningful baseline."""
    if not first:
        return None
    return _round((last - first) / first * 100.0)


def aggregate_training_data(
    items: Iterable[dict],
    now: datetime,
    months: int = 6,
) -> dict:
    """Filter records to the last `months` and aggregate into a summary dict.

    `now` is passed in (not read from the clock) so callers/tests are deterministic.
    """
    cutoff = now - timedelta(days=round(_DAYS_PER_MONTH * months))

    # Normalize + filter to the window. Parse timestamps once.
    records: list[dict] = []
    for it in items:
        try:
            ts = datetime.fromisoformat(it["created_at"])
        except (KeyError, ValueError):
            continue
        if ts < cutoff:
            continue
        records.append(
            {
                "name": it["name"],
                "weight": float(it["weight"]),
                "reps": int(it["reps"]),
                "ts": ts,
                "day": ts.date(),
            }
        )

    period = {
        "start": cutoff.date().isoformat(),
        "end": now.date().isoformat(),
        "months": months,
    }

    if not records:
        return {"period": period, "summary": _empty_summary(), "monthly": [], "exercises": []}

    # ---- overall summary ----
    workout_days = {r["day"] for r in records}
    total_volume = sum(r["weight"] * r["reps"] for r in records)
    window_weeks = max((now.date() - cutoff.date()).days / 7.0, 1.0)
    summary = {
        "total_sets": len(records),
        "total_reps": sum(r["reps"] for r in records),
        "total_volume": _round(total_volume),
        "distinct_exercises": len({r["name"] for r in records}),
        "workout_days": len(workout_days),
        "avg_workouts_per_week": _round(len(workout_days) / window_weeks),
        "first_workout": min(workout_days).isoformat(),
        "last_workout": max(workout_days).isoformat(),
    }

    # ---- monthly time series (volume/frequency trend for the LLM) ----
    by_month: dict[str, list[dict]] = defaultdict(list)
    for r in records:
        by_month[r["ts"].strftime("%Y-%m")].append(r)
    monthly = [
        {
            "month": month,
            "sets": len(rs),
            "total_reps": sum(r["reps"] for r in rs),
            "total_volume": _round(sum(r["weight"] * r["reps"] for r in rs)),
            "workout_days": len({r["day"] for r in rs}),
        }
        for month, rs in sorted(by_month.items())
    ]

    # ---- per-exercise progression ----
    by_exercise: dict[str, list[dict]] = defaultdict(list)
    for r in records:
        by_exercise[r["name"]].append(r)

    exercises = []
    for name, rs in by_exercise.items():
        # Group this exercise's sets by day to derive first/last session trends.
        days: dict = defaultdict(list)
        for r in rs:
            days[r["day"]].append(r)
        ordered_days = sorted(days)
        first_day, last_day = ordered_days[0], ordered_days[-1]
        first_session = _session_stats(days[first_day])
        last_session = _session_stats(days[last_day])

        total_volume_ex = sum(r["weight"] * r["reps"] for r in rs)
        exercises.append(
            {
                "name": name,
                "sessions": len(ordered_days),
                "sets": len(rs),
                "total_reps": sum(r["reps"] for r in rs),
                "total_volume": _round(total_volume_ex),
                "max_weight": max(r["weight"] for r in rs),
                "best_est_1rm": _round(max(estimated_1rm(r["weight"], r["reps"]) for r in rs)),
                "avg_weight": _round(sum(r["weight"] for r in rs) / len(rs)),
                "avg_reps": _round(sum(r["reps"] for r in rs) / len(rs)),
                "first_session": {"date": first_day.isoformat(), **first_session},
                "last_session": {"date": last_day.isoformat(), **last_session},
                "est_1rm_trend_pct": _pct_change(first_session["est_1rm"], last_session["est_1rm"]),
                "volume_trend_pct": _pct_change(first_session["volume"], last_session["volume"]),
            }
        )

    # Most-trained first — helps the LLM prioritize.
    exercises.sort(key=lambda e: e["total_volume"], reverse=True)

    return {"period": period, "summary": summary, "monthly": monthly, "exercises": exercises}


def _empty_summary() -> dict:
    return {
        "total_sets": 0,
        "total_reps": 0,
        "total_volume": 0,
        "distinct_exercises": 0,
        "workout_days": 0,
        "avg_workouts_per_week": 0,
        "first_workout": None,
        "last_workout": None,
    }
