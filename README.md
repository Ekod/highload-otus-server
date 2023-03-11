# highload-otus-server

**Запуск сервиса**
1. Убедиться, что установлен Docker-Compose
2. Из корневой директории проекта запустить 
   ```shell 
   make up
   ```

**Миграции**

```shell
# Установка библиотеки миграции
brew install golang-migrate
# Создание миграции
migrate create -ext sql -dir zarf/migrations <migration_name_change_me>
# Применение миграции (на локальную БД)
make migrate
```

**Остановка сервиса**
1. Из корневой директории проекта запустить
 ```shell 
   make down
 ```