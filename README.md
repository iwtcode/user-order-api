# user-order-api

## Описание

REST API для управления пользователями и их заказами на Go с использованием PostgreSQL, GORM, Gin и JWT-авторизации.

---

## Основные эндпоинты

> ✅ **Защищено** — требуется JWT-токен в заголовке запроса `Authorization: Bearer <token>`

| Метод  | Эндпоинт                      | Описание                              | Авторизация                           |
|--------|-------------------------------|---------------------------------------|---------------------------------------|
| POST   | /auth/login                   | Авторизация                           | <div align="center">🔓</div>          |
| POST   | /users                        | Создание пользователя                 | <div align="center">🔓</div>          |
| GET    | /users                        | Получение списка пользователей        | <div align="center">🔒</div>          |
| GET    | /users/{id}                   | Получение пользователя по ID          | <div align="center">🔒</div>          |
| PUT    | /users/{id}                   | Обновление пользователя               | <div align="center">🔒</div>          |
| DELETE | /users/{id}                   | Удаление пользователя                 | <div align="center">🔒</div>          |
| POST   | /users/{user_id}/orders       | Создание заказа для пользователя      | <div align="center">🔒</div>          |
| GET    | /users/{user_id}/orders       | Получение списка заказов пользователя | <div align="center">🔒</div>          |

Полная документация — [Swagger UI](http://localhost:8080/swagger/index.html)

## Быстрый старт

### 1. Запуск через Docker Compose (рекомендуется)

#### 1.1. Установите и запустите Docker

> Требуется установленный [Docker](https://www.docker.com/) и [Docker Compose](https://docs.docker.com/compose/)

#### 1.2. Создайте файл **.env** и настройте переменные окружения (Опционально)

```ini
SERVER_PORT=8080            # Порт, на котором запускается сервер
DB_HOST=localhost           # Адрес сервера базы данных PostgreSQL
DB_PORT=5432                # Порт базы данных PostgreSQL
DB_USER=postgres            # Имя пользователя для подключения к БД
DB_PASSWORD=db-password     # Пароль пользователя для подключения к БД
DB_NAME=user_order_api      # Имя базы данных
DB_SSLMODE=disable          # Режим SSL для подключения к БД
JWT_SECRET=your-secret-key  # Секретный ключ для подписи JWT
JWT_EXPIRATION=24h          # Время жизни токена (например, 24h)
GIN_MODE=debug              # Режим работы Gin (debug/release)
LOG_FILE=logs/app.log       # Путь к файлу логов
```

#### 1.3. Запустите построение образа Docker

```sh
docker compose up -d --build
```

#### 1.4. Запустите приложение

```
docker compose up
```

#### Дополнительно. Смена переменных окружения (очистка базы данных и создание новой) / При непредвиденных ошибках

```
docker compose down
docker volume rm user-order-api_db-data
docker compose up -d --build
```

---

### 2. Локальный запуск без Docker

#### 2.1. Установите Golang и PostgreSQL
> Требуется установленный [Golang (1.24.2)](https://go.dev/dl/) и [PostgreSQL](https://www.postgresql.org/download/)

#### 2.2. Создайте файл **.env** и настройте переменные окружения (Обязательно)

```ini
SERVER_PORT=8080            # Порт, на котором запускается сервер
DB_HOST=localhost           # Адрес сервера базы данных PostgreSQL
DB_PORT=5432                # Порт базы данных PostgreSQL
DB_USER=postgres            # Имя пользователя для подключения к БД
DB_PASSWORD=db-password     # Пароль пользователя для подключения к БД
DB_NAME=user_order_api      # Имя базы данных
DB_SSLMODE=disable          # Режим SSL для подключения к БД
JWT_SECRET=your-secret-key  # Секретный ключ для подписи JWT
JWT_EXPIRATION=24h          # Время жизни токена (например, 24h)
GIN_MODE=debug              # Режим работы Gin (debug/release)
LOG_FILE=logs/app.log       # Путь к файлу логов
```

#### 2.3. Запустите установщик (через bash/git bash)

```
./setup.sh
```

#### 2.4. Запустите приложение

```
./server.exe
```
или
```
go run cmd/main.go
```

---

### Структура проекта

```
user-order-api/
├── cmd/                    # Точка входа в приложение
├── internal/               # Внутренние пакеты приложения
│   ├── config/             # Конфигурация приложения
│   ├── handlers/           # HTTP-обработчики (контроллеры)
│   ├── middleware/         # Промежуточные обработчики (обёртки)
│   ├── models/             # Описания моделей данных (структуры для БД)
│   ├── repository/         # Слой доступа к данным (репозитории)
│   ├── services/           # Бизнес-логика (сервисы)
│   ├── utils/              # Вспомогательные функции и утилиты
│   └── tests/              # Тесты
├── migrations/             # SQL-миграции
├── docs/                   # Документация API
├── go.mod                  # Файл зависимостей Go-модулей
├── go.sum                  # Контрольные суммы зависимостей
├── LICENSE                 # Лицензия проекта
├── README.md               # Основная документация и инструкции по запуску
├── setup.sh                # Установщик миграций и приложения
├── Dockerfile              # Dockerfile для сборки контейнера приложения
└── docker-compose.yml      # Docker Compose для запуска приложения и БД
```