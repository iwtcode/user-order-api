# user-order-api

## Описание

REST API для управления пользователями и их заказами на Go с использованием PostgreSQL, GORM, Gin и JWT-авторизации.

---

## Быстрый старт

> **Документация API доступна после запуска приложения по адресу:**
> [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

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

---

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
├── internal/                    # Внутренние пакеты приложения
│   ├── config/                  # Конфигурация приложения
│   ├── handlers/                # HTTP-обработчики (контроллеры)
│   ├── middleware/              # Промежуточные обработчики (обёртки)
│   ├── models/                  # Описания моделей данных (структуры для БД)
│   ├── repository/              # Слой доступа к данным (репозитории)
│   ├── services/                # Бизнес-логика (сервисы)
│   └── utils/                   # Вспомогательные функции и утилиты
├── migrations/                  # SQL-миграции
├── docs/                        # Документация API
├── tests/                       # Тесты
├── go.mod                       # Файл зависимостей Go-модулей
├── go.sum                       # Контрольные суммы зависимостей
├── LICENSE                      # Лицензия проекта
├── README.md                    # Основная документация и инструкции по запуску
├── setup.sh                     # Установщик миграций и приложения
├── Dockerfile                   # Dockerfile для сборки контейнера приложения
└── docker-compose.yml           # Docker Compose для запуска приложения и БД
```