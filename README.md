# GophKeeper

GophKeeper представляет собой клиент-серверную систему, 
позволяющую пользователю надёжно и безопасно хранить логины, 
пароли, бинарные данные и прочую приватную информацию.


# HTTP API сервиса

## Регистрация
- POST `/api/v1/auth/signup`
    ```json
    {
      "login": "example_login",
      "password": "example_pwd"
    }
    ```
  - 201 - пользователь зарегестрирован
  - 409 - пользователь уже существует
  - 500 - внутренняя ошибка сервера
  - 
## Логин пользователя
- POST `/api/v1/auth/signin`
    ```json
    {
      "login": "example_login",
      "password": "example_pwd"
    }
    ```
  - 201 - пользователь аутентифицирован
    ```json
    {
      "token": "1q2w3e4r5t"
    }
    ```
  - 404 - пользователь не найден
  - 409 - пользователь уже существует
  - 500 - внутренняя ошибка сервера

## 
- DELETE `/api/v1/auth/signout/`
    ```
    HTTP Header 
         Token: "1q2w3e4r5t"
    ```
- POST `/api/v1/secret/create`
  ```
    HTTP Header 
         Token: "1q2w3e4r5t"
  ```
  ```json
  {
    // ???
  }
  ```
- PUT `/api/v1/secret/update`
  ```
    HTTP Header 
         Token: "1q2w3e4r5t"
  ```
  ```json
  {
    // ???
  }
  ```
- GET `/api/v1/secret/{secretID}`

- DELETE `/api/v1/secret/{secretID}`