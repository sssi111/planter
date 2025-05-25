# Planter - Plant Care Backend

Planter is a backend service for a mobile application that helps users care for their plants, flowers, and bushes. It provides features such as plant recommendations based on user preferences, care instructions, reminders for watering and fertilizing, and a knowledge base with verified information.

## Features

- **User Authentication**: Register and login with email/password or Google authentication
- **Plant Database**: Comprehensive database of plants with care instructions
- **Plant Recommendations**: AI-powered plant recommendations based on user preferences
- **Care Reminders**: Notifications for watering, fertilizing, and repotting
- **Favorites**: Save favorite plants for quick access
- **User Plants**: Track plants owned by users with watering history
- **Shop Integration**: Browse plants available in shops

## Tech Stack

- **Backend**: Go (Golang)
- **Database**: PostgreSQL
- **AI Integration**: Yandex GPT for plant recommendations
- **Authentication**: JWT-based authentication
- **API Documentation**: OpenAPI 3.0
- **Containerization**: Docker

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL (or use the Docker Compose setup)
- Yandex GPT API key

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```
# Server
PORT=8080

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=planter
DB_SSLMODE=disable

# Authentication
JWT_SECRET=your-secret-key
TOKEN_DURATION=24

# Yandex GPT
YANDEX_GPT_API_KEY=your-yandex-gpt-api-key
YANDEX_GPT_MODEL=yandexgpt
```

### Running with Docker

1. Build and start the containers:

```bash
docker-compose up -d
```

2. The API will be available at `http://localhost:8080`

### Running Locally

1. Install dependencies:

```bash
go mod download
```

2. Run the application:

```bash
go run cmd/api/main.go
```

## API Documentation

The API is documented using OpenAPI 3.0. You can find the documentation in the `docs/openapi.yaml` file.

## Database Schema

The database schema is defined in the `scripts/schema.sql` file. It includes tables for:

- Users
- Plants
- Care Instructions
- User Plants
- User Favorite Plants
- Shops
- Shop Plants
- Special Offers
- Plant Questionnaires
- Plant Recommendations

## Project Structure

```
.
├── cmd/
│   └── api/              # Application entry point
├── docs/
│   └── openapi.yaml      # API documentation
├── internal/
│   ├── api/              # API handlers
│   ├── auth/             # Authentication
│   ├── config/           # Configuration
│   ├── db/               # Database connection
│   ├── middleware/       # Middleware
│   ├── models/           # Data models
│   ├── repository/       # Data access layer
│   │   └── impl/         # Repository implementations
│   ├── services/         # Business logic
│   └── utils/            # Utilities
├── pkg/
│   ├── logger/           # Logging
│   └── validator/        # Validation
├── scripts/
│   └── schema.sql        # Database schema
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## Testing

There are several ways to run the tests:

### Using the test script

We've provided a script that runs all tests and generates a coverage report:

```bash
# Make the script executable if needed
chmod +x scripts/run_tests.sh

# Run the tests
./scripts/run_tests.sh
```

This will:
1. Run all tests with verbose output
2. Generate a coverage report
3. Create an HTML coverage report at `coverage.html`

### Running tests manually

You can also run tests manually with the Go test command:

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/services

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Running tests in Docker

You can run tests inside the Docker container:

```bash
docker-compose run --rm api go test ./...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.