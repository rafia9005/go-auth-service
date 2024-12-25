## API Documentation

### Authentication Service

### 1. Register User

**Endpoint**: `/register`

**Method**: `POST`

**Description**: Registers a new user.

**Request Body**:
```json
{
    "first_name": "John",
    "last_name": "Doe",
    "email": "johndoe@example.com",
    "password": "password123"
}
```

**Responses**:

- **200 OK**:
  ```json
  {
      "status": true,
      "message": "Registration successful! Please check your email for the verification code"
  }
  ```
- **400 Bad Request**:
  ```json
  {
      "message": "Invalid request payload"
  }
  ```
- **409 Conflict**:
  ```json
  {
      "message": "Email already in use"
  }
  ```
- **500 Internal Server Error**:
  ```json
  {
      "message": "Failed to register user"
  }
  ```

### 2. Login User

**Endpoint**: `/login`

**Method**: `POST`

**Description**: Authenticates a user and returns access and refresh tokens.

**Request Body**:
```json
{
    "email": "johndoe@example.com",
    "password": "password123"
}
```

**Responses**:

- **200 OK**:
  ```json
  {
      "status": true,
      "access_token": "access-token-example",
      "refresh_token": "refresh-token-example"
  }
  ```
- **400 Bad Request**:
  ```json
  {
      "message": "Invalid request payload"
  }
  ```
- **401 Unauthorized**:
  ```json
  {
      "message": "Invalid email or password"
  }
  ```
- **403 Forbidden**:
  ```json
  {
      "message": "Account not verified. Please check your email for verification instructions."
  }
  ```
- **500 Internal Server Error**:
  ```json
  {
      "message": "Error generating access token"
  }
  ```

### 3. Authenticate with Google

**Endpoint**: `/auth/google`

**Method**: `GET`

**Description**: Redirects to Google authentication URL.

**Query Parameters**:
- `from` (optional): The URL to redirect back to after authentication.

**Responses**:

- **302 Found**: Redirects to Google authentication URL.

### 4. Google Authentication Callback

**Endpoint**: `/auth/google/callback`

**Method**: `GET`

**Description**: Handles Google authentication callback and returns access and refresh tokens.

**Query Parameters**:
- `code`: Authorization code from Google.

**Responses**:

- **200 OK**:
  ```json
  {
      "status": "success",
      "message": "User authenticated successfully",
      "access_token": "access-token-example",
      "refresh_token": "refresh-token-example",
      "data": {
          "user": {
              "ID": 1,
              "Name": "John Doe",
              "FirstName": "John",
              "LastName": "Doe",
              "Email": "johndoe@example.com",
              "CreatedAt": "2024-12-25T03:21:53Z",
              "UpdatedAt": "2024-12-25T03:21:53Z"
          }
      }
  }
  ```
- **400 Bad Request**:
  ```json
  {
      "status": "error",
      "message": "Authorization code is missing"
  }
  ```
- **401 Unauthorized**:
  ```json
  {
      "status": "error",
      "message": "Failed to exchange authorization code for token"
  }
  ```
- **500 Internal Server Error**:
  ```json
  {
      "status": "error",
      "message": "Failed to get user info: error message"
  }
  ```

### 5. Authenticate with GitHub

**Endpoint**: `/auth/github`

**Method**: `GET`

**Description**: Redirects to GitHub authentication URL.

**Query Parameters**:
- `from` (optional): The URL to redirect back to after authentication.

**Responses**:

- **302 Found**: Redirects to GitHub authentication URL.

### 6. GitHub Authentication Callback

**Endpoint**: `/auth/github/callback`

**Method**: `GET`

**Description**: Handles GitHub authentication callback and returns access and refresh tokens.

**Query Parameters**:
- `code`: Authorization code from GitHub.

**Responses**:

- **200 OK**:
  ```json
  {
      "status": "success",
      "message": "User authenticated successfully",
      "access_token": "access-token-example",
      "refresh_token": "refresh-token-example",
      "data": {
          "user": {
              "ID": 1,
              "Name": "John Doe",
              "FirstName": "John",
              "LastName": "Doe",
              "Email": "johndoe@example.com",
              "CreatedAt": "2024-12-25T03:21:53Z",
              "UpdatedAt": "2024-12-25T03:21:53Z"
          }
      }
  }
  ```
- **400 Bad Request**:
  ```json
  {
      "status": "error",
      "message": "Authorization code is missing"
  }
  ```
- **401 Unauthorized**:
  ```json
  {
      "status": "error",
      "message": "Failed to exchange authorization code for token"
  }
  ```
- **500 Internal Server Error**:
  ```json
  {
      "status": "error",
      "message": "Failed to get user info: error message"
  }
  ```

### 7. Verify Token

**Endpoint**: `/verify-token`

**Method**: `GET`

**Description**: Verifies the provided token.

**Headers**:
- `x-token`: Access token to be verified.

**Responses**:

- **200 OK**:
  ```json
  {
      "status": true,
      "message": "Token is valid",
      "claims": "claims object"
  }
  ```
- **401 Unauthorized**:
  ```json
  {
      "status": "false",
      "message": "Unauthorized: Token is missing"
  }
  ```
  ```json
  {
      "status": "false",
      "message": "Invalid Token"
  }
  ```

### 8. Refresh Token

**Endpoint**: `/refresh-token`

**Method**: `POST`

**Description**: Refreshes the access token using the provided refresh token.

**Request Body**:
```json
{
    "refresh_token": "refresh-token-example"
}
```

**Responses**:

- **200 OK**:
  ```json
  {
      "access_token": "new-access-token-example",
      "refresh_token": "new-refresh-token-example"
  }
  ```
- **400 Bad Request**:
  ```json
  {
      "message": "Invalid request payload"
  }
  ```
- **401 Unauthorized**:
  ```json
  {
      "message": "Invalid or expired refresh token"
  }
  ```
- **500 Internal Server Error**:
  ```json
  {
      "message": "Failed to generate access token"
  }
  ```
