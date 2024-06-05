# 🧊 Тестовое задание


# 🚀 Быстрый старт

1. Склонируй репозиторий

```bash 
  git clone https://github.com/limona77/Auth-JWT
```
2. Запусти программу
```bash 
  go run cmd/main.go
 ``` 

## 🌐 API документация

### Регистрация

- `POST /auth/register`
- **Запрос**:
  ```json
  {
    "email": "user1@gmail.com",
    "password": "12345"
  }
  ```

- **Ответ**:
  ```json
  {
    "user": {
        "ID": 3,
        "Email": "user1@gmail.com",
        "Password": ""
    },
    "refreshToken": "token",
    "accessToken": "token"
  }
  ```


### Логин
- `POST /auth/login`
- **Запрос**:
  ```json
  {
    "email": "user1@gmail.com",
    "password": "12345"
  }
  
  ```

- **Ответ**:
  ```json
  {
    "user": {
        "ID": 3,
        "Email": "user1@gmail.com",
        "Password": ""
    },
    "refreshToken": "token",
    "accessToken": "token"
  }
   "birthdays": [
        "ДР у пользователя user2@gmail.com: 05-06-2024"
    ]
  ```

### Обновление токенов
- `GET auth/refresh`
- **Запрос**: body не нужно
- **Cookie**: должен храниться refreshToken
- **Ответ**:
  ```json
  {
    "user": {
        "ID": 3,
        "Email": "user1@gmail.com",
        "Password": ""
    },
    "refreshToken": "token",
    "accessToken": "token"
  }
  ```
### Выход из аккаунта
- `GET /auth/logout`
- **Запрос**: body не нужно
- **Cookie**: должен храниться refreshToken
- **Ответ**:
  ```json
  {"userID": 3}
  ```

### Получить свои данные
- `GET /me`
- **Запрос**: body не нужно
- **Authorization=Bearer accessToken**

- **Ответ**:
  ```json
  {
    "user": {
        "ID": 3,
        "Email": "user1@gmail.com",
        "Password": ""
    }
  }
  ```
### Подписка на пользователя
- `POST subscribe?id={id пользователя}`
- **Запрос**: body не нужно
- **Authorization=Bearer accessToken**

- **Ответ**:
  ```json
  {
    "subscribe": {
        "UserID": 3,
        "SubscribedToId": 4
    }
  }
  ```
### Отписка от пользователя
- `DELETE unsubscribe?id={id пользователя}`
- **Запрос**: body не нужно
- **Authorization=Bearer accessToken**

- **Ответ**:
  ```json
  {
    "subscribe": {
        "UserID": 3,
        "SubscribedToId": 4
    }
  }
  ```