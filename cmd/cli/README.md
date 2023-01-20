# About

cli клиент для работы с сервисом хранения секретов

# конфигурация cli клиента

```shell
export GK_API="https://apiserver.example.com"
#  Если необходимо изменить пути лог файла и хранилища токена сессии
export GK_LOG="path/to/file"
export GK_TOKEN="path/to/file"
```

По умолчанию значения:
```shell
GK_API="http://127.0.0.1:8080"
GK_LOG="~/.gk_cli_R2D2"
GK_TOKEN="~/.gk_cli_logs"
```

# Getting start

cli клиент работает на основе [cobra](https://github.com/spf13/cobra)
Для каждого уровня команд есть help. 

```shell
go run cmd/cli/main.go
Goph Keeper command line interface

Usage:
  gkc [command]

Available Commands:
  help        Help about any command
  secret      secret service
  signin      LogIn
  signout     Logout user
  signup      Create new account

Flags:
  -h, --help   help for gkc

Use "gkc [command] --help" for more information about a command.

```

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

**Создать новый секрет**
```shell
gk secret new --secret='{"login": "qwe", "password": "asd"}' --name="Пароль от сейфа" --description=""
```

**Получить список секретов**
```shell
gk secret list
```

**Получить секрет по id**
```shell
gk secret get --id=3
```

**Обновить секрет по id**
```shell
gk secret update --id=3 --secret='{"login": "qwe1", "password": "asd123"}' --name="" --description=""
```
В случае если не переданы один или несколько флагов:
- `--secret`
- `--name`
- `--description`
значения будут подставлены из тех, что сейчас определены для секрета

**Удалить секрет по id**
```shell
gk secret delete --id=3
```