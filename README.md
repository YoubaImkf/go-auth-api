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
```

### Configuration

Edit `configs/config.yaml` to configure database and SMTP settings.

### Running the Application

1. Build and run the application using Docker Compose:

    ```sh
    docker-compose up --build -d
    ```

2. Access the API at `http://localhost:8080`.

### Swagger Documentation

Access the Swagger UI at `http://localhost:8080/swagger/index.html`.

## 📚 API Endpoints

### Auth

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Authenticate a user
- `POST /api/v1/auth/forgot-password` - Request a password reset
- `POST /api/v1/auth/reset-password` - Reset the user's password
- `POST /api/v1/auth/logout` - Logout a user (protected)
- `GET /api/v1/auth/me` - Get user profile (protected)

### Health

- `GET /api/v1/health` - Check the health of the service
