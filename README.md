# TrueAuth API

The TrueAuth API is a standalone authentication gRPC and REST API service that provides various endpoints for user authentication.

#### Database Overview at DB Docs

- ##### Preview [dbdocs.io/sirjager/trueauth](https://dbdocs.io/sirjager/trueauth)
- ##### Password [github.com/sirjager/trueauth](https://dbdocs.io/sirjager/trueauth)

## Service Description

The TrueAuth service provides the following functionalities:

- **Welcome**: Returns a welcome message.
- **Health**: Checks the health of the API.
- **Register**: Registers a new user.
- **Login**: Authenticates a user and generates an access token.
- **Verify**: Requests and verifies an email.
- **Logout**: Logs out active sessions.
- **Refresh**: Refreshes an access token.
- **Recovery**: Generates a password recovery code.
- **Update**: Updates user details.
- **Delete**: Requests user deletion.
- **AllowIP**: Allows an IP address.
- **User**: Retrieves information about a user.

## Endpoints

### Welcome

- **Description**: Returns a welcome message.
- **Summary**: Welcome Message
- **Endpoint**: `GET /`
- **Tags**: System
- **Example**:

  ```ts
  GET http://localhost:4421
  content-type: application/json
    {
        "message": "Welcome to TrueAuth Api"
    }
  ```

### Health

- **Description**: Checks the health of the API.
- **Summary**: API Health
- **Endpoint**: `GET /v1/health`
- **Tags**: System
- **Example**:
  ```ts
  GET http://localhost:4421/v1/health
  content-type: application/json
    {
        "status": "UP",
        "uptime": "2066.368459083s",
        "started": "2023-05-02T10:25:15.394742686Z",
        "timestamp": "2023-05-02T10:59:41.763201267Z"
    }
  ```

### Register

- **Description**: Registers a new user.
- **Summary**: Register User
- **Endpoint**: `POST /v1/register`
- **Tags**: Auth
- **Example**:

  ```ts
  POST http://localhost:4421/v1/register
  content-type: application/json
    {
        "email": "email@gmail.com",
        "username": "johndoe",
        "password": "password"
    }

  RESPONSE
  {
    "user": {
        "id": "2fb0eea2-9b1f-4103-a206-ec50547193e7",
        "email": "email@gmail.com",
        "username": "johndoe",
        "created_at": "2023-05-02T11:04:57.044452Z",
        "updated_at": "2023-05-02T11:04:57.044452Z"
    }
  }
  ```

### Login

- **Description**: Authenticates a user and generates an access token.
- **Summary**: Login User
- **Endpoint**: `POST /v1/login`
- **Tags**: Auth
- **Example**:

  ```ts
  POST http://localhost:4421/v1/login
  content-type: application/json
    {
        "identity": "email@gmail.com",
        "password": "password"
    }

  RESPONSE
    "user": {
        "id": "99488f4d-b8ab-4ea6-8722-b0f1fe2051d3",
        "email": "email@gmail.com",
        "username": "johndoe",
        "created_at": "2023-05-02T17:37:45.970720Z",
        "updated_at": "2023-05-02T17:37:45.970720Z"
    },
    "session_id": "048cfdb8-b92f-449f-9bf5-6de4c8ca7b64",
    "access_token": "v2.local.L2t4A6f96-Qd0gXqBVEZxfcu0BJdJE-Ywwp3w0unicGoCn2DisFjV_AE06Ch-z6pqmE_4kO7VjBsHipLFpoSNnGacwpDoqTsI9w081uJ1lkkpNlwwoWQIjwEVDJ1DgxKqkU9PhDTB8HzFS1R67CFSi-FlCzvU4_sSz_ze9ygojK91gsv_jtAMiLlvDn6s7n3bk5pqiXNowlylVjHkHmiwvQSO7TGjPmAyx6rOT122EG1HUED__uTEJD2iVumjcCTqEwgGi8kyb90bR3RZvdcB2XoupaoOe_wbvcInxxA4pIU8ufUBV-W5EqOERNqeg.bnVsbA",
    "refresh_token": "v2.local.FM9i9wmpWYe9Ag9JEjTpPpgMldyj5lKLGVCshKfjUvL3apJfbff6ehsNLjPIYO-Y4h0B4e8XWzoNWgTS7CzH9i0TLE8wR0oJXdmttf8bdlAkX3NhcUrTxLbcJPcQxQH0peNOWLBSxtw-o7hS214WpfVfXb49hgzb2yqB9I71Kcs_MaTHmNxbAzOrA2TAqn8r3U9jGwqPalMob4K67ZGBL3OOcgzsQIAenSKCwmArrIRTptdr81nEi4esIWozhKTWY2GcRcIX2G6s9qVqL6rrYBCIxoCq4RtVFN--saGRRcba7r9CbnohpBxDAv4QQA.bnVsbA",
    "access_token_expires_at": "2023-05-05T17:37:56.677812136Z",
    "refresh_token_expires_at": "2023-05-08T17:37:56.678011913Z"
  ```

### Verify

- **Description**: Requests and verifies an email.
- **Summary**: Request and verify email
- **Endpoint**: `POST /v1/verify`
- **Tags**: Auth
- **Example: When email is not verified**

  ```ts
  POST http://localhost:4421/v1/verify
  content-type: application/json
  Authorization: bearer access_token

  RESPONSE
    {
        "message": "email verification code sent to your email email@gmail.com"
    }
  ```

- **Example: When already have verification code**

  ```ts
  POST http://localhost:4421/v1/verify?code=123456
  content-type: application/json
  Authorization: bearer access_token

  RESPONSE
    {
        "message": "email email@gmail.com successfully verified"
    }
  ```

### Logout

- **Description**: Logs out active sessions.
- **Summary**: Logout
- **Endpoint**: `POST /v1/logout`
- **Tags**: Auth
- **Example: Logout all sessions**:

  ```ts
  POST http://localhost:4421/v1/logout
  content-type: application/json
  Authorization: bearer access_token

  RESPONSE
  {
    "message": "sessions deleted"
  }
  ```

- **Example: Logout specific sessions**:

  ```ts
  POST http://localhost:4421/v1/logout?session=f9d4f78f-8fb2-454b-acd4-520a0c1ff8d6,0f3feae3-5369-407b-846b-880ae783976b
  content-type: application/json
  Authorization: bearer access_token

  RESPONSE
  {
    "message": "sessions deleted"
  }
  ```

### Refresh

- **Description**: Refreshes an access token.
- **Summary**: Refresh access token
- **Endpoint**: `POST /v1/refresh`
- **Tags**: Auth
- **Example**:

  ```ts
  POST http://localhost:4421/v1/refresh
  content-type: application/json
  {
    "refresh_token": "refresh-token",
  }

  RESPONSE
  {
    "access_token": "new-access-token",
    "access_token_expires_at": "2023-05-05T11:34:27.572306237Z"
  }
  ```

### Recovery

- **Description**: Generates a password recovery code.
- **Summary**: Password recovery
- **Endpoint**: `POST /v1/recovery`
- **Tags**: Auth
- **Example: When requesting for password reset**:

  ```ts
  POST http://localhost:4421/v1/recovery
  content-type: application/json
  {
    "email": "email@gmail.com",
  }

  RESPONSE
  {
    "message": "password recovery code has been sent to your email email@gmail.com"
  }
  ```

- **Example: When already have password reset code**:

  ```ts
  POST http://localhost:4421/v1/recovery
  content-type: application/json
  {
    "email": "email@gmail.com",
    "code": "123456"
  }

  RESPONSE
  {
    "message": "your new password has been sent to your email"
  }
  ```

### Update

- **Description**: Updates user details.
- **Summary**: Update user
- **Endpoint**: `PATCH /v1/users`
- **Tags**: Users
- **Example**:

  ```ts
  PATCH http://localhost:4421/v1/users
  content-type: application/json
  Authorization: bearer access_token
  {
    "username": "johndoe",
    "password": "newpassword",
    "firstname": "john",
    "lastname": "doe"
  }

  RESPONSE
  {
    "user": {
        "id": "2fb0eea2-9b1f-4103-a206-ec50547193e7",
        "email": "email@gmail.com",
        "username": "johndoe",
        "firstname": "john",
        "lastname": "doe",
        "verified": true,
        "created_at": "2023-05-02T11:04:57.044452Z",
        "updated_at": "2023-05-02T11:49:05.106080Z"
    }
  }
  ```

### Delete

- **Description**: Requests user deletion.
- **Summary**: Delete user
- **Endpoint**: `DELETE /v1/users`
- **Tags**: Auth
- **Example: When email is verified**:

  ```ts
  DELETE http://localhost:4421/v1/users
  content-type: application/json
  Authorization: bearer access_token

  RESPONSE
  {
    "message": "deletion code has been sent to your email"
  }
  ```

- **Example: When already have deletion code**:

  ```ts
  DELETE http://localhost:4421/v1/users?code=123456
  content-type: application/json
  Authorization: bearer access_token

  RESPONSE
  {
    "message": "user account has been deleted"
  }
  ```

### AllowIP

- **Description**: Allows an IP address.
- **Summary**: Allow IP address
- **Endpoint**: `GET /v1/allowip`
- **Tags**: Auth
- **Example: When requesting from unknown ip address**:

  ```ts
  GET http://localhost:4421/v1/allowip
  content-type: application/json
  Authorization: bearer access_token

  RESPONSE
  {
    "message": "code has been sent to your email to allow requests from this ip address"
  }
  ```

- **Example: When already have allow ip code**:

  ```ts
  GET http://localhost:4421/v1/allowip?code=123456
  content-type: application/json
  Authorization: bearer access_token

  RESPONSE
  {
    "message": "your ip address has been successfully added to whitelist"
  }
  ```

### User

- **Description**: Retrieves information about a user.
- **Summary**: Returns requested user
- **Endpoint**: `GET /v1/users/{identity}`
- **Tags**: Users
- **Example**:

  ```ts
  GET http://localhost:4421/v1/users/2fb0eea2-9b1f-4103-a206-ec50547193e7
  GET http://localhost:4421/v1/users/email@gmail.com
  GET http://localhost:4421/v1/users/johndoe
  content-type: application/json
  Authorization: bearer access_token

  RESPONSE
  {
    "user": {
        "id": "2fb0eea2-9b1f-4103-a206-ec50547193e7",
        "email": "email@gmail.com",
        "username": "johndoe",
        "firstname": "john",
        "lastname": "doe",
        "verified": true,
        "created_at": "2023-05-02T11:04:57.044452Z",
        "updated_at": "2023-05-02T11:52:03.616376Z"
    }
  }
  ```

## Additional Information

For more information about the TrueAuth API and its usage, please refer to the [TrueAuth GitHub repository](https://github.com/sirjager/trueauth). You can find detailed documentation, examples, and further resources there.
