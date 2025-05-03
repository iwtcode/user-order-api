# user-order-api

## Описание

REST API для управления пользователями и их заказами на Go с использованием PostgreSQL, GORM, Gin и JWT-авторизации.

---

## Быстрый старт

### 1. Запуск через Docker Compose (рекомендуется)

#### 1.1. Установите и запустите Docker

> Требуется установленный [Docker](https://www.docker.com/) и [Docker Compose](https://docs.docker.com/compose/)

#### 1.2. Создайте файл **.env** и настройте переменные окружения (Опционально)

```ini
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-postgres-password
DB_NAME=userorderapi
DB_SSLMODE=disable
JWT_SECRET=your-secret-key
```

#### 1.3. Запустите построение образа Docker

```sh
docker compose up -d --build
```

#### 1.4. Запустите приложение

```
docker compose up
```

#### Дополнительно. Смена переменных окружения (очистка базы данных и создание новой)

```
docker compose down
docker volume rm user-order-api_db-data
docker compose up -d --build
```

### 2. Локальный запуск без Docker

#### 2.1. Установите Golang и PostgreSQL
> Требуется установленный [Golang (1.24.2)](https://go.dev/dl/) и [PostgreSQL](https://www.postgresql.org/download/)

#### 2.2. Создайте файл **.env** и настройте переменные окружения (Обязательно)

```ini
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-postgres-password
DB_NAME=userorderapi
DB_SSLMODE=disable
JWT_SECRET=your-secret-key
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
├── cmd/                         # Точка входа в приложение
│   └── main.go                  # Главный файл запуска сервера
├── internal/                    # Внутренние пакеты приложения
│   ├── config/                  # Конфигурация приложения
│   │   └── config.go            # Загрузка переменных окружения
│   ├── handlers/                # HTTP-обработчики (контроллеры)
│   │   ├── user_handler.go      # Обработчики для пользователей
│   │   ├── order_handler.go     # Обработчики для заказов
│   │   └── auth_handler.go      # Обработчик для авторизации
│   ├── models/                  # Описания моделей данных (структуры для БД)
│   │   ├── user.go              # Модель пользователя
│   │   └── order.go             # Модель заказа
│   ├── repository/              # Слой доступа к данным (репозитории)
│   │   ├── user_repo.go         # Методы работы с пользователями в БД
│   │   └── order_repo.go        # Методы работы с заказами в БД
│   ├── services/                # Бизнес-логика (сервисы)
│   │   ├── user_service.go      # Логика управления пользователями
│   │   ├── order_service.go     # Логика управления заказами
│   │   └── auth_service.go      # Логика авторизации и работы с токенами
│   ├── middleware/              # Промежуточные обработчики (обёртки)
│   │   └── auth_middleware.go   # JWT-аутентификация для защищённых маршрутов
│   └── utils/                   # Вспомогательные функции и утилиты
│       ├── jwt.go               # Генерация и валидация JWT-токенов
│       └── password.go          # Хеширование и проверка паролей
├── migrations/                  # SQL-миграции
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_orders_table.up.sql
│   └── 000002_create_orders_table.down.sql
├── docs/                        # Документация API
├── go.mod                       # Файл зависимостей Go-модулей
├── go.sum                       # Контрольные суммы зависимостей
├── .env                         # Пример файла переменных окружения
├── .gitignore                   # Файлы и папки, игнорируемые Git
├── .dockerignore                # Файлы и папки, игнорируемые Docker
├── Dockerfile                   # Dockerfile для сборки контейнера приложения
├── docker-compose.yml           # Docker Compose для запуска приложения и БД
├── setup.sh                     # Установщик миграций и приложения
├── README.md                    # Основная документация и инструкции по запуску
└── LICENSE                      # Лицензия проекта
```