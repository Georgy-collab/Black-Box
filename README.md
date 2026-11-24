# User Service

REST API для управления пользователями на FastAPI. Данные хранятся в памяти (без БД).

## Возможности

- создание, чтение, обновление и удаление пользователей
- валидация email, имени и возраста
- проверка уникальности email
- фильтрация списка по минимальному возрасту
- health-check

## Стек

- Python 3.10+
- FastAPI
- Pydantic v2
- Uvicorn

## Установка и запуск

```bash
pip install -r requirements.txt
uvicorn app:app --reload
```

Сервис будет доступен по адресу: http://127.0.0.1:8000

Интерактивная документация:

- Swagger UI: http://127.0.0.1:8000/docs
- ReDoc: http://127.0.0.1:8000/redoc
- OpenAPI JSON: http://127.0.0.1:8000/openapi.json

Статические файлы спецификации в репозитории:

- [`openapi.yaml`](openapi.yaml)
- [`openapi.json`](openapi.json)

Обновить файлы после изменения API:

```bash
python export_openapi.py
```

## API

| Метод  | Путь             | Описание                                      |
|--------|------------------|-----------------------------------------------|
| GET    | `/health`        | Проверка работоспособности                    |
| POST   | `/users`         | Создать пользователя                          |
| GET    | `/users`         | Список пользователей (`?min_age=` опционально)|
| GET    | `/users/{id}`    | Получить пользователя по id                   |
| PUT    | `/users/{id}`    | Обновить пользователя                         |
| DELETE | `/users/{id}`    | Удалить пользователя                          |

### Модель пользователя

```json
{
  "name": "Alice",
  "email": "alice@example.com",
  "age": 25
}
```

Правила валидации:

- `name` — непустая строка (1–100 символов)
- `email` — корректный email, сохраняется в нижнем регистре
- `age` — целое число от 1 до 149
- email должен быть уникальным

### Примеры

Создать пользователя:

```bash
curl -X POST http://127.0.0.1:8000/users \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"Alice\", \"email\": \"alice@example.com\", \"age\": 25}"
```

Список с фильтром:

```bash
curl "http://127.0.0.1:8000/users?min_age=18"
```

## Структура проекта

```
├── app.py              # FastAPI-приложение и эндпоинты
├── validators.py       # Валидация email и возраста
├── openapi.yaml        # OpenAPI 3 спецификация
├── openapi.json        # OpenAPI 3 (JSON)
├── export_openapi.py   # Экспорт OpenAPI из приложения
├── requirements.txt    # Зависимости
└── README.md
```

## Примечание

Хранилище `users_db` — словарь в памяти. После перезапуска сервиса данные сбрасываются.
