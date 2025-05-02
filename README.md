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
docker compose up --build
```

### Дополнительно

#### Повторные запуски

```
docker compose up
```

#### Смена переменных окружения (очистка базы данных и создание новой)

```
docker compose down
docker volume rm user-order-api_db-data
docker compose up --build
```

### 2. Локальный запуск без Docker

#### 2.1. Установите Golang
> Требуется установленный [Golang (1.24.2)](https://go.dev/dl/)

#### 2.2. Создайте файл **.env** и настройте переменные окружения

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