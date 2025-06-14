# Receipt Backend API

Go backend API for receipt management application with Google OAuth authentication and MySQL database.

## Features

- Google OAuth 2.0 authentication
- JWT token-based authorization
- MySQL database with GORM
- RESTful API with Gin framework
- UUID v7 for user IDs
- Soft delete support
- CORS enabled

## Setup

1. Copy environment variables:
```bash
cp .env.example .env
```

2. Configure your `.env` file with:
   - Database connection details
   - JWT secret key
   - Google OAuth credentials

3. Install dependencies:
```bash
go mod tidy
```

4. Run the application:
```bash
go run cmd/app/main.go
```

## API Endpoints

### Authentication

- `POST /api/v1/auth/google` - Authenticate with Google OAuth
- `GET /api/v1/auth/validate` - Validate JWT token

### Protected Routes

- `GET /api/v1/profile` - Get user profile (requires authentication)

### Health Check

- `GET /health` - Health check endpoint

## Database Schema

### Users Table
- `id` - Binary UUID v7 (primary key)
- `created_at` - Creation timestamp
- `updated_at` - Update timestamp  
- `deleted_at` - Soft delete timestamp (nullable)
- `google_id` - Google OAuth user ID (nullable, unique)

## Authentication Flow

1. Frontend receives Google OAuth access token
2. Frontend sends token to `/api/v1/auth/google`
3. Backend validates token with Google API
4. Backend checks if user exists in database
5. If user doesn't exist, creates new user record
6. Backend generates JWT token with user ID
7. Backend returns JWT token to frontend
8. Frontend uses JWT token for subsequent API calls

## Deployment

### Simple deploy from source:
```bash
gcloud run deploy receipt-backend --region=europe-west3 --allow-unauthenticated --source .
```

### Docker build and deploy:
```bash
GOOS=linux GOARCH=amd64 go build -o bin/app cmd/app/main.go 
docker buildx build --platform linux/amd64 -t gcr.io/money-advice-462707/receipt-backend -f bin/Dockerfile .
docker push gcr.io/money-advice-462707/receipt-backend
```