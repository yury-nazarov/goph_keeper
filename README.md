# GophKeeper

GophKeeper представляет собой клиент-серверную систему, 
позволяющую пользователю надёжно и безопасно хранить логины, 
пароли, бинарные данные и прочую приватную информацию.

# Схема приложения

**В данном разделе отражается только реализованный на данный момент функционал**

## Dev

Подходит для тестового развертывания/разработки, не рекомендуется для продакшн.

![](docs/img/app.png)

## Prod

Рекомендуемая схема:
- Терминация TLS на nginx, а не в app т.к. шифрование довольно дорогая для CPU операция.
  Совместно с бизнес логикой приложения и обильным логированием может привести к довольно интересным последствиям.
  В данном случае nginx может быть как отдельный инстанс так и ingress k8s если это позволяет дизайн кластера.

![](docs/img/app_nginx.png)

# Локальный запуск приложения

Запустить Psql из docker-compose в директории `deployments`
```shell
docker-compose up
```

Запустить сервер с ключами 
```shell
go run cmd/gophkeeper/main.go -a 127.0.0.1:8080 -d "host=localhost port=5432 user=gop_keeper_dev password=gop_keeper_dev dbname=gop_keeper_dev sslmode=disable connect_timeout=5"
```
или переменными окружения
```shell
export RUN_ADDRESS="127.0.0.1:8080"
export DATABASE_URI="host=localhost port=5432 user=gop_keeper_dev password=gop_keeper_dev dbname=gop_keeper_dev sslmode=disable connect_timeout=5"
go run cmd/gophkeeper/main.go
```

# HTTP API сервиса

## Регистрация пользователя
**Request**

POST `/api/v1/auth/signup`
 ```json
{
  "login": "example_login",
  "password": "example_pwd"
}
```

**Response**  

- 201 - пользователь зарегестрирован
  ```
  HTTP Authorization: {{Token}}
  ```
- 400 - не верный формат зароса
- 409 - пользователь уже существует
- 500 - внутренняя ошибка сервера

## Вход пользователя
**Request**

POST `/api/v1/auth/signin`
```json
{
  "login": "example_login",
  "password": "example_pwd"
}
```

**Response**
- 201 - пользователь аутентифицирован
  ```
  HTTP Authorization: {{Token}}
  ```
- 404 - пользователь не найден
- 409 - пользователь уже существует
- 500 - внутренняя ошибка сервера

## Выход пользователя

**Request**

DELETE `/api/v1/auth/signout/`
```
HTTP Authorizations: {{Token}}
```

**Response**

- 200 - пользовательская сессия удалена

## Создать секрет

**Request**

POST `/api/v1/secret/create`
```
HTTP Authorization: {{Token}}
```

**Response**

TODO

## Обновить секрет

**Request**

PUT `/api/v1/secret/update`

**Response**

## Получить секрет

**Request**

GET `/api/v1/secret/{secretID}`

**Response**

TODO

## Список секретов

**Request**

GET `/api/v1/secret/list`

**Response**

TODO

##  Удалить секрет

**Request**

DELETE `/api/v1/secret/{secretID}`

**Response**