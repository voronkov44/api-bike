# API-Bike 🚲

Backend-сервис для управления пользователями и продуктами для проекта bike. Написан на Go с использованием Docker и PostgreSQL.

## 🛠 Технологии
- Go 1.24 (чистая архитектура)
- PostgreSQL (хранение данных)
- GORM (ORM для работы с БД)
- Docker (контейнеризация приложения)
- Docker Compose (инструмент для запуска нескольких контейнеров через единый конфигурационный файл)
- Swagger (документация API)


## Установка и запуск

1. Клонируем репозиторий:
```bash
git clone https://github.com/voronkov44/api-bike.git
```
2. Заходим в корень проекта:
```bash
cd api-bike/
```

3. Создание enviroments-file:
```bash
nano .env
```
Пример файла:
```
POSTGRES_USER=postgres 
POSTGRES_PASSWORD=postgres 
POSTGRES_DB=bike 
POSTGRES_PORT=5432 
DSN=host=postgres user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} port=${POSTGRES_PORT} sslmode=disable 
SECRET=1
```
#### Пояснение:
DSN — строка подключения к базе данных PostgreSQL.

Важно:
Параметр host=postgres менять не нужно, так как это имя контейнера базы данных внутри Docker-Compose сети.
Контейнеры в одной сети Docker взаимодействуют между собой по именам сервисов из docker-compose.yml.

SECRET — секретный ключ для генерации JWT-токенов.
Установите здесь любой надёжный ключ для защиты авторизации в API.

4. Запуск
*Требуется установка [docker](https://www.docker.com/products/docker-desktop/), если не установлен, смотрите [зависимости.]()*
```
docker-compose up -d
```

Сервер будет доступен на http://localhost:8081

### Общие команды для работы с контейнерами

```
# Просмотр запущенных контейнеров
docker ps

#Просмотр всех контейнеров, включая остановленные
docker ps -a

# Запуск в фоновом режиме
docker-compose up -d --build

# Остановка всех сервисов
docker-compose down

# Остановка контейнера
docker stop <container_id>

# Удлание контейнера
docker rm <container_id>

# Удаление образа
docker rmi <image_id>

# Просмотр логов контейнера
docker logs -f <container_name>

# Вход в контейнер
docker exec -it <container_name> sh - sh

docker exec -it <container_name> /bin/bash - bash

# Очистка системы
docker system prune
```

## Зависимости
### Установка docker
Установка пакета [Docker Engine](https://docs.docker.com/engine/install/)

