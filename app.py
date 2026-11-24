from typing import Optional

from fastapi import FastAPI, HTTPException, Path, Query, status
from pydantic import BaseModel, Field, field_validator

from validators import is_valid_age, is_valid_email

tags_metadata = [
    {
        "name": "health",
        "description": "Проверка состояния сервиса.",
    },
    {
        "name": "users",
        "description": "CRUD-операции с пользователями.",
    },
]

app = FastAPI(
    title="User Service",
    description=(
        "REST API для управления пользователями.\n\n"
        "Данные хранятся в памяти процесса и сбрасываются при перезапуске."
    ),
    version="1.0.0",
    contact={"name": "Black-Box", "url": "https://github.com/Georgy-collab/Black-Box"},
    openapi_tags=tags_metadata,
    docs_url="/docs",
    redoc_url="/redoc",
    openapi_url="/openapi.json",
)

users_db: dict[int, dict] = {}
next_id = 1


class UserCreate(BaseModel):
    name: str = Field(
        min_length=1,
        max_length=100,
        description="Имя пользователя",
        examples=["Alice"],
    )
    email: str = Field(
        description="Email пользователя (уникальный)",
        examples=["alice@example.com"],
    )
    age: int = Field(
        ge=1,
        le=149,
        description="Возраст пользователя",
        examples=[25],
    )

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
    name: Optional[str] = Field(
        default=None,
        min_length=1,
        max_length=100,
        description="Новое имя пользователя",
        examples=["Alicia"],
    )
    email: Optional[str] = Field(
        default=None,
        description="Новый email пользователя",
        examples=["alicia@example.com"],
    )
    age: Optional[int] = Field(
        default=None,
        ge=1,
        le=149,
        description="Новый возраст пользователя",
        examples=[26],
    )

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
    id: int = Field(description="Идентификатор пользователя", examples=[1])
    name: str = Field(description="Имя пользователя", examples=["Alice"])
    email: str = Field(description="Email пользователя", examples=["alice@example.com"])
    age: int = Field(description="Возраст пользователя", examples=[25])


class HealthResponse(BaseModel):
    status: str = Field(description="Статус сервиса", examples=["ok"])


class ErrorResponse(BaseModel):
    detail: str = Field(description="Описание ошибки", examples=["User not found"])


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


@app.get(
    "/health",
    tags=["health"],
    summary="Проверка работоспособности",
    description="Возвращает статус сервиса.",
    response_model=HealthResponse,
    responses={200: {"description": "Сервис доступен"}},
)
def health():
    return {"status": "ok"}


@app.post(
    "/users",
    tags=["users"],
    summary="Создать пользователя",
    description="Создаёт нового пользователя. Email должен быть уникальным.",
    status_code=status.HTTP_201_CREATED,
    response_model=UserResponse,
    responses={
        201: {"description": "Пользователь создан"},
        400: {"model": ErrorResponse, "description": "Email уже занят или некорректные данные"},
        422: {"description": "Ошибка валидации входных данных"},
    },
)
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


@app.get(
    "/users",
    tags=["users"],
    summary="Список пользователей",
    description="Возвращает всех пользователей. Можно отфильтровать по минимальному возрасту.",
    response_model=list[UserResponse],
    responses={200: {"description": "Список пользователей"}},
)
def list_users(
    min_age: Optional[int] = Query(
        default=None,
        ge=0,
        description="Минимальный возраст для фильтрации",
        examples=[18],
    ),
):
    result = list(users_db.values())
    if min_age is not None:
        result = [user for user in result if user["age"] >= min_age]
    return result


@app.get(
    "/users/{user_id}",
    tags=["users"],
    summary="Получить пользователя",
    description="Возвращает пользователя по идентификатору.",
    response_model=UserResponse,
    responses={
        200: {"description": "Пользователь найден"},
        404: {"model": ErrorResponse, "description": "Пользователь не найден"},
    },
)
def get_user(
    user_id: int = Path(description="Идентификатор пользователя", examples=[1], ge=1),
):
    return _get_user_or_404(user_id)


@app.put(
    "/users/{user_id}",
    tags=["users"],
    summary="Обновить пользователя",
    description="Частично обновляет данные пользователя. Можно передать любые поля.",
    response_model=UserResponse,
    responses={
        200: {"description": "Пользователь обновлён"},
        400: {"model": ErrorResponse, "description": "Email уже занят или некорректные данные"},
        404: {"model": ErrorResponse, "description": "Пользователь не найден"},
        422: {"description": "Ошибка валидации входных данных"},
    },
)
def update_user(
    user: UserUpdate,
    user_id: int = Path(description="Идентификатор пользователя", examples=[1], ge=1),
):
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


@app.delete(
    "/users/{user_id}",
    tags=["users"],
    summary="Удалить пользователя",
    description="Удаляет пользователя по идентификатору.",
    status_code=status.HTTP_204_NO_CONTENT,
    responses={
        204: {"description": "Пользователь удалён"},
        404: {"model": ErrorResponse, "description": "Пользователь не найден"},
    },
)
def delete_user(
    user_id: int = Path(description="Идентификатор пользователя", examples=[1], ge=1),
):
    _get_user_or_404(user_id)
    del users_db[user_id]
    return None
