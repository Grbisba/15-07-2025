# Техническое задание

____

## Сводное HTTP API


### Система работы с пользователем должна предоставлять следующие HTTP-хендлеры:

Система работы с пользователем доступна по пути среднего уровня `/users`.

- **POST** `/users/register` — регистрация пользователя.
- **POST** `/users/login` — аутентификация пользователя.
- **POST** `/users/refresh-token` — перевыпуск пары авторизационных токенов.

### Система работы с задачами должна предоставлять следующие HTTP-хендлеры:

Система работы с задачами доступна по пути среднего уровня `/tasks`.

- **POST** `/tasks/create` - создание задачи.
- **POST** `/tasks/add-file/:task_id` - добавление файлов в задачу.
- **GET** `/tasks/status/:task_id` - получение статуса задачи.
- **GET** `/tasks/download/:task_id` - загрузка ".zip" файла.

### Описание методов системы для работы с пользователем.

- #### Регистрация пользователя

##### Хендлер: **POST** `/users/registration`.

Формат запроса:

```json lines
POST /users/registration HTTP/1.1
Content-Type application/json
...
{
  "login": "<login>",
  "password": "<password>"
} 
```

Формат ответа:

```json lines
HTTP/1.1
Content-Type application/json
...

{
  "access_token": "<access_token>",
  "refresh_token": "<refresh_token>"
}
```

- #### Аутентификация пользователя

##### Хендлер: **POST** `/users/login`.

Формат запроса:

```json lines
POST /users/login HTTP/1.1
Content-Type application/json
...

{
  "login": "<login>",
  "password": "<password>"
}
```

Формат ответа:

```json lines
HTTP/1.1
Content-Type application/json
...

{
  "access_token": "<access_token>",
  "refresh_token": "<refresh_token>"
}
```


- #### Перевыпуск пары авторизационных токенов.

##### Хендлер: **POST** `/users/refresh-token`.

Формат запроса:

```json lines
POST /users/refresh-token HTTP/1.1
Content-Type application/json
Authorization Bearer <token>
...
```

Формат ответа:

```json lines
HTTP/1.1
Content-Type application/json
...

{
  "access_token": "<access_token>",
  "refresh_token": "<refresh_token>"
}
```

### Описание методов системы для работы с задачами.

- #### Создание задачи

##### Хендлер: **POST** `/tasks/create`.

Формат запроса:

```json lines
POST /tasks/create HTTP/1.1
Content-Type application/json
Authorization Bearer <token>
...
```

Формат ответа:

```json lines
HTTP/1.1
Content-Type application/json
...
{  
  "id": "<id>",
  "status": "<status>"
}
```

- #### Добавление файлов в задачу

##### Хендлер: **POST** `/tasks/add-file/:task_id`.

Формат запроса:

```json lines
POST /tasks/add-file/:task_id HTTP/1.1
Content-Type application/json
Authorization Bearer <token>
...
{
  "urls": ["<url1>", "<url2>", ...]
}
```

Формат ответа:

```json lines
HTTP/1.1
Content-Type application/json
...
{  
  "id": "<id>",
  "status": "<status>"
}
```

- #### Получение статуса задачи

##### Хендлер: **GET** `/tasks/status/:task_id`.

Формат запроса:

```json lines
GET /tasks/status/:task_id HTTP/1.1
Content-Type application/json
Authorization Bearer <token>
...
```

Формат ответа:

```json lines
HTTP/1.1
Content-Type application/json
...
{
  "user_id": "<user_id>",
  "status": "<status>",
  "url": "<url>",
  "uploaded_files": ["<url1>", "<url2>", ...]
}
```

- #### Загрузка ".zip" файла

##### Хендлер: **GET** `/tasks/download/:task_id`.

Формат запроса:

```json lines
GET /tasks/download/:task_id HTTP/1.1
Content-Type application/json
...
```

Формат ответа:

```json lines
HTTP/1.1
Content-Type application/zip
Content-Disposition attachment filename="<filename>.zip"
...
<raw zip-file>
```