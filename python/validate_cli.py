#!/usr/bin/env python3
"""CLI bridge: read JSON fields from stdin, write validation result to stdout."""

import json
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))

from validators import validate_fields  # noqa: E402


def main() -> int:
    try:
        payload = json.load(sys.stdin)
    except json.JSONDecodeError:
        json.dump({"ok": False, "detail": "Invalid JSON"}, sys.stdout)
        return 1

    if not isinstance(payload, dict):
        json.dump({"ok": False, "detail": "Payload must be an object"}, sys.stdout)
        return 1

    try:
        normalized = validate_fields(payload)
    except ValueError as exc:
        json.dump({"ok": False, "detail": str(exc)}, sys.stdout)
        return 1

    json.dump({"ok": True, "data": normalized}, sys.stdout)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
