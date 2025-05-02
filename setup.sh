#!/bin/bash

# Загружаем переменные из .env
if [ -f .env ]; then
  export $(grep -v '^#' .env | sed -e 's/^export //g' | xargs)
else
  echo ".env file not found!"
  exit 1
fi

# Проверка наличия обязательных переменных
: "${DB_USER?Need to set DB_USER}"
: "${DB_PASSWORD?Need to set DB_PASSWORD}"
: "${DB_HOST?Need to set DB_HOST}"
: "${DB_PORT?Need to set DB_PORT}"
: "${DB_NAME?Need to set DB_NAME}"
: "${DB_SSLMODE:=disable}"
: "${MIGRATIONS_PATH=./migrations}"

echo "Checking if database '$DB_NAME' exists..."

# Пытаемся подключиться к стандартной БД 'postgres' для проверки
if PGPASSWORD="$DB_PASSWORD" psql -U "$DB_USER" -h "$DB_HOST" -p "$DB_PORT" -lqt | cut -d \| -f 1 | grep -qw "$DB_NAME"; then
    echo "Database '$DB_NAME' already exists."
else
    echo "Database '$DB_NAME' does not exist. Creating..."
    if PGPASSWORD="$DB_PASSWORD" psql -U "$DB_USER" -h "$DB_HOST" -p "$DB_PORT" -d postgres -c "CREATE DATABASE \"$DB_NAME\""; then
        echo "Database '$DB_NAME' created successfully."
    else
        echo "Failed to create database '$DB_NAME'."
        exit 1
    fi
fi

# Устанавливаем утилиту migrate, если ее нет или обновляем до последней версии
echo "Installing/updating migrate CLI tool..."
if go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; then
  echo "migrate CLI installed/updated successfully."
else
  echo "Failed to install migrate CLI. Make sure Go is installed and configured correctly."
  exit 1
fi


echo "Running database migrations..."
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

# Проверяем, доступна ли команда migrate
if ! command -v migrate &> /dev/null
then
    echo "'migrate' command could not be found."
    echo "Please ensure $(go env GOPATH)/bin or $HOME/go/bin is in your PATH."
    exit 1
fi

# Применяем миграции
if migrate -path "${MIGRATIONS_PATH}" -database "${DATABASE_URL}" up; then
  echo "Migrations applied successfully."
else
  echo "Failed to apply migrations."
  exit 1
fi

echo "Downloading Go modules..."
if go mod download; then
  echo "Go modules downloaded successfully."
else
  echo "Failed to download Go modules."
  exit 1
fi

echo "Building Go application..."
if go build -o server.exe cmd/main.go; then
  echo "Build successful. Binary: ./server"
else
  echo "Build failed."
  exit 1
fi

echo "Setup complete."
exit 0