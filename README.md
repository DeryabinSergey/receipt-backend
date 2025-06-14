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
- Optimized for Google Cloud Run

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

## Cloud Run Deployment

The application is optimized for Google Cloud Run with:
- Intelligent migration system that only runs when needed
- Graceful shutdown handling
- Optimized logging for Cloud Run
- Proper timeout configurations

### Deploy to Cloud Run

1. **Simple deploy from source:**
```bash
gcloud run deploy receipt-backend \
  --region=europe-west3 \
  --allow-unauthenticated \
  --source . \
  --set-env-vars="GIN_MODE=release"
```

2. **Docker build and deploy:**
```bash
# Build for linux/amd64
docker buildx build --platform linux/amd64 -t gcr.io/YOUR_PROJECT_ID/receipt-backend .

# Push to Container Registry
docker push gcr.io/YOUR_PROJECT_ID/receipt-backend

# Deploy to Cloud Run
gcloud run deploy receipt-backend \
  --image gcr.io/YOUR_PROJECT_ID/receipt-backend \
  --region=europe-west3 \
  --allow-unauthenticated
```

### Environment Variables for Cloud Run

Set these environment variables in Cloud Run:
- `DB_HOST` - Your Cloud SQL instance connection
- `DB_PORT` - Database port (usually 3306)
- `DB_USER` - Database username
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `JWT_SECRET` - Your JWT secret key
- `GIN_MODE` - Set to "release" for production

## Migration Strategy

The application uses intelligent migrations that:
- Check if migration is actually needed before running
- Only migrate when table structure changes
- Skip unnecessary operations on each cold start
- Log migration status for debugging

This approach minimizes startup time in Cloud Run while ensuring database consistency.