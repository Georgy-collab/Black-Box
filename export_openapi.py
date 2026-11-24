"""Export OpenAPI schema to openapi.json and openapi.yaml."""

from pathlib import Path

import yaml

from app import app

ROOT = Path(__file__).resolve().parent
schema = app.openapi()

json_path = ROOT / "openapi.json"
yaml_path = ROOT / "openapi.yaml"

json_path.write_text(
    __import__("json").dumps(schema, ensure_ascii=False, indent=2) + "\n",
    encoding="utf-8",
)
yaml_path.write_text(
    yaml.dump(schema, allow_unicode=True, sort_keys=False),
    encoding="utf-8",
)

print(f"Wrote {json_path.name} and {yaml_path.name}")
