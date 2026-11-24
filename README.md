# User Service

Гибридный REST API: **HTTP-сервис на Go**, **валидация на Python**.

Данные хранятся в памяти (без БД).

## Архитектура

```
Клиент → Go (net/http) → Python validate_cli.py → Go store → ответ
```

- Go принимает запросы, хранит пользователей, отдаёт ответы
- Python нормализует и проверяет `name`, `email`, `age`
- Связь: Go запускает `python/validate_cli.py` и обменивается JSON через stdin/stdout

## Возможности

- CRUD пользователей
- валидация email, имени и возраста (Python)
- уникальность email
- фильтр списка по `min_age`
- health-check
- Swagger UI и OpenAPI

## Требования

- Go 1.22+
- Python 3.10+ (стандартная библиотека, без pip-зависимостей)

## Запуск

Из корня репозитория:

```bash
go run ./cmd/server
```

### Docker

```bash
docker build -t user-service .
docker run --rm -p 8000:8000 --name user-service user-service
```

Сервис: http://127.0.0.1:8000

Документация:

- Swagger UI: http://127.0.0.1:8000/docs
- OpenAPI: http://127.0.0.1:8000/openapi.json

Опционально:

- `PORT` — порт сервиса
- `PYTHON_BIN` — путь к Python (по умолчанию `python`)
- `PYTHON_VALIDATOR` — путь к `validate_cli.py`

## API

| Метод  | Путь             | Описание                                      |
|--------|------------------|-----------------------------------------------|
| GET    | `/health`        | Проверка работоспособности                    |
| POST   | `/users`         | Создать пользователя                          |
| GET    | `/users`         | Список пользователей (`?min_age=` опционально)|
| GET    | `/users/{id}`    | Получить пользователя по id                   |
| PUT    | `/users/{id}`    | Обновить пользователя                         |
| DELETE | `/users/{id}`    | Удалить пользователя                          |

### Модель

```json
{
  "name": "Alice",
  "email": "alice@example.com",
  "age": 25
}
```

Правила:

- `name` — непустая строка (1–100 символов)
- `email` — корректный email, в нижнем регистре
- `age` — целое число от 1 до 149
- email уникален

## Структура

```
├── cmd/server/              # точка входа Go
├── internal/
│   ├── handlers/            # HTTP handlers
│   ├── models/              # DTO
│   ├── pyvalidate/          # вызов Python-валидатора
│   ├── server/              # роутинг и docs
│   └── store/               # in-memory хранилище
├── python/
│   ├── validators.py        # правила валидации
│   └── validate_cli.py      # CLI-мост для Go
├── openapi.yaml
├── openapi.json
├── Dockerfile
├── .dockerignore
├── go.mod
└── README.md
```

## Тесты

```bash
go test ./...
```

Запускать из корня репозитория, чтобы Go нашёл `python/validate_cli.py`.
