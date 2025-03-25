FROM golang:1.24-alpine

WORKDIR /app

# Installe swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Générer la documentation Swagger
RUN swag init -d ./cmd/api/ --pdl 3

RUN go build -o main cmd/api/main.go

EXPOSE 8080

CMD ["./main"]