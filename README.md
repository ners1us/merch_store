# merch_store
Приложение, работающее с магазином мерча.

## Инструменты
- Go
- Gin
- SQL
- Docker
- PostgreSQL
- Testcontainers

## Запуск приложения
```bash
docker compose up -d --build
```

## Остановка приложения
```bash
docker compose down
```

## Просмотр логов приложения
```bash
docker logs merch_store-app-1
```

## Очистка базы данных
```bash
docker volume rm merch_store_postgres_data
```

## Данные для авторизации в БД
- PostgreSQL, порт - 5432:
    - username: admin
    - password: password

## Примечания
- Спецификация API находится в файле ```schema.json```
- Работу endpoint'ов рекомендуется проверять в Postman.
