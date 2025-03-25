# Authentication API

A JWT Authentication Service API built with Go, Gin, and Gorm.

## 🚀 Features

- User registration and login
- JWT-based authentication
- Password reset via email
- Token blacklisting for logout
- Swagger documentation

## 🛠️ Setup

### Prerequisites

- Docker
- Docker Compose

### Environment Variables

Create a `.env` file with the following content:

```env
APP_ENVIRONMENT=development
APP_HOST=localhost:8080

# Database configuration for PostgreSQL
DATABASE_HOST=db # 'db' with docker-compose, 'localhost' if in local
DATABASE_PORT=5432
DATABASE_USER=root
DATABASE_PASSWORD=root
DATABASE_NAME=go-auth-db

# SMTP configuration
SMTP_HOST=smtp # 'smtp' with docker-compose, 'localhost' if in local
SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM=no-reply@example.com

# JWT secret
JWT_SECRET={secret}
```

### Configuration

Edit `configs/config.yaml` to configure jwt and group settings.

### Running the Application

1. Build and run the application using Docker Compose:

    ```sh
    docker-compose up --build -d
    ```

2. Access the API at `http://localhost:8080`.

3. Access the SMTP at `http://localhost:1080`.

### Running Locally

To run the application locally without Docker:

1. Ensure you updated the **docker-compose.yml** and **config.yml** files as comments suggest it !

2. Generate the swagger documentation :

    ```sh
    swag inti -g cmd/api/main.go
    ```
    or 
    ```sh
    swag init -d ./cmd/api/ --pdl 3
    ```

3. Ensure PostgreSQL and MailDev are running:

    ```sh
    docker-compose up db smtp -d
    ```

4. Run the application:

    ```sh
    go run cmd/api/main.go
    ```

### Swagger Documentation

Access the Swagger UI at `http://localhost:8080/swagger/index.html`.

## 📚 API Endpoints

### Auth

- `POST /{UUID}/register` - Register a new user
- `POST /{UUID}/login` - Authenticate a user
- `POST /{UUID}/forgot-password` - Request a password reset
- `POST /{UUID}/reset-password` - Reset the user's password
- `POST /{UUID}/logout` - Logout a user (protected)
- `GET /{UUID}/me` - Get user profile (protected)

### Health

- `GET /{UUID}/health` - Check the health of the service

### USer

- `GET /{UUID}/users` - Check the health of the service
- `DELETE /{UUID}/remove-users` - Check the health of the service


## 🧪 Running Tests

1. To run tests, use the following command:

    ```sh
    go test ./...
    ```