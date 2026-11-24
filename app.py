from typing import Optional

from fastapi import FastAPI, HTTPException, Query, status
from pydantic import BaseModel, Field, field_validator

from validators import is_valid_age, is_valid_email

app = FastAPI(title="User Service")

users_db: dict[int, dict] = {}
next_id = 1


class UserCreate(BaseModel):
    name: str = Field(min_length=1, max_length=100)
    email: str
    age: int = Field(ge=1, le=149)

    @field_validator("email")
    @classmethod
    def validate_email(cls, value: str) -> str:
        email = value.strip().lower()
        if not is_valid_email(email):
            raise ValueError("Invalid email")
        return email

    @field_validator("name")
    @classmethod
    def validate_name(cls, value: str) -> str:
        name = value.strip()
        if not name:
            raise ValueError("Name cannot be empty")
        return name


class UserUpdate(BaseModel):
    name: Optional[str] = Field(default=None, min_length=1, max_length=100)
    email: Optional[str] = None
    age: Optional[int] = Field(default=None, ge=1, le=149)

    @field_validator("email")
    @classmethod
    def validate_email(cls, value: Optional[str]) -> Optional[str]:
        if value is None:
            return value
        email = value.strip().lower()
        if not is_valid_email(email):
            raise ValueError("Invalid email")
        return email

    @field_validator("name")
    @classmethod
    def validate_name(cls, value: Optional[str]) -> Optional[str]:
        if value is None:
            return value
        name = value.strip()
        if not name:
            raise ValueError("Name cannot be empty")
        return name


class UserResponse(BaseModel):
    id: int
    name: str
    email: str
    age: int


def _get_user_or_404(user_id: int) -> dict:
    user = users_db.get(user_id)
    if user is None:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="User not found")
    return user


def _ensure_unique_email(email: str, exclude_user_id: Optional[int] = None) -> None:
    for user_id, user in users_db.items():
        if exclude_user_id is not None and user_id == exclude_user_id:
            continue
        if user["email"] == email:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail="Email already exists",
            )


@app.get("/health")
def health():
    return {"status": "ok"}


@app.post("/users", status_code=status.HTTP_201_CREATED, response_model=UserResponse)
def create_user(user: UserCreate):
    global next_id

    if not is_valid_age(user.age):
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="Invalid age")

    _ensure_unique_email(user.email)

    user_id = next_id
    users_db[user_id] = {
        "id": user_id,
        "name": user.name,
        "email": user.email,
        "age": user.age,
    }
    next_id += 1
    return users_db[user_id]


@app.get("/users", response_model=list[UserResponse])
def list_users(min_age: Optional[int] = Query(default=None, ge=0)):
    result = list(users_db.values())
    if min_age is not None:
        result = [user for user in result if user["age"] >= min_age]
    return result


@app.get("/users/{user_id}", response_model=UserResponse)
def get_user(user_id: int):
    return _get_user_or_404(user_id)


@app.put("/users/{user_id}", response_model=UserResponse)
def update_user(user_id: int, user: UserUpdate):
    stored = _get_user_or_404(user_id)

    if user.name is not None:
        stored["name"] = user.name

    if user.email is not None:
        _ensure_unique_email(user.email, exclude_user_id=user_id)
        stored["email"] = user.email

    if user.age is not None:
        if not is_valid_age(user.age):
            raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="Invalid age")
        stored["age"] = user.age

    return stored


@app.delete("/users/{user_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_user(user_id: int):
    _get_user_or_404(user_id)
    del users_db[user_id]
    return None
