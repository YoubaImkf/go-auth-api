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
POSTGRES_USER={root}
POSTGRES_PASSWORD={root}
POSTGRES_DB={db}

JWT_SECRET={your_jwt_secret_here}

DATABASE_HOST=localhost
DATABASE_PORT=5432

SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=
```

### Configuration

Edit `configs/config.yaml` to configure database and SMTP settings.

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

2. Ensure PostgreSQL and MailDev are running:

    ```sh
    docker-compose up db smtp -d
    ```

3. Run the application:

    ```sh
    go run cmd/api/main.go
    ```

### Swagger Documentation

Access the Swagger UI at `http://localhost:8080/swagger/index.html`.

## 📚 API Endpoints

### Auth

- `POST /{UUID}/auth/register` - Register a new user
- `POST /{UUID}/auth/login` - Authenticate a user
- `POST /{UUID}/auth/forgot-password` - Request a password reset
- `POST /{UUID}/auth/reset-password` - Reset the user's password
- `POST /{UUID}/auth/logout` - Logout a user (protected)
- `GET /{UUID}/auth/me` - Get user profile (protected)

### Health

- `GET /{UUID}/health` - Check the health of the service

### USer

- `GET /{UUID}/users` - Check the health of the service


## 🧪 Running Tests

1. To run tests, use the following command:

    ```sh
    go test ./...
    ```