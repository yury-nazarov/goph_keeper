# About

cli клиент для работы с сервисом хранения секретов

# Getting start
## Auth
Регистрация нового пользователя
```shell
gk signup --login=admin --password=123
```

Вход пользователя
```shell
gk signin --login=admin --password=123
```

Выход пользователя
```shell
gk signout
```

## Secrets

Создать новый секрет
```shell
gk secret new --secret='{"login": "qwe", "password": "asd"}'
```

Получить список секретов
```shell
gk secret list
```

Получить секрет по id
```shell
gk secret get --id=3
```

Обновить секрет по id
```shell
gk secret update --id=3 --secret='{"login": "qwe1", "password": "asd123"}'
```

Удалить секрет по id
```shell
gk secret delete --id=3
```