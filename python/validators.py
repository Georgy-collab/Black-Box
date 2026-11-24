import re

EMAIL_PATTERN = re.compile(r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$")


def is_valid_email(email: str) -> bool:
    return bool(EMAIL_PATTERN.fullmatch(email))


def is_valid_age(age: int) -> bool:
    return 1 <= age <= 149


def normalize_name(name: str) -> str:
    name = name.strip()
    if not name:
        raise ValueError("Name cannot be empty")
    if len(name) > 100:
        raise ValueError("Name is too long")
    return name


def normalize_email(email: str) -> str:
    email = email.strip().lower()
    if not is_valid_email(email):
        raise ValueError("Invalid email")
    return email


def validate_fields(payload: dict) -> dict:
    """Validate and normalize provided user fields.

    Accepts any subset of: name, email, age.
    Returns normalized values for the same keys.
    """
    result = {}

    if "name" in payload:
        if not isinstance(payload["name"], str):
            raise ValueError("Name must be a string")
        result["name"] = normalize_name(payload["name"])

    if "email" in payload:
        if not isinstance(payload["email"], str):
            raise ValueError("Email must be a string")
        result["email"] = normalize_email(payload["email"])

    if "age" in payload:
        age = payload["age"]
        if isinstance(age, bool) or not isinstance(age, int):
            raise ValueError("Invalid age")
        if not is_valid_age(age):
            raise ValueError("Invalid age")
        result["age"] = age

    return result
