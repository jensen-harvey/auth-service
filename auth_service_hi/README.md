# Authentication Service for Herb Immortal

This is a simple authentication service built with Go that provides user registration, login, and session management for Customers, Admins, Healers, and Vendors.

## Features

- **User Registration**: Supports different user roles (Customer, Admin, Healer, Vendor)
- **User Authentication**: Email/password login with JWT token generation
- **Session Management**: Session validation and management
- **Role-Based Authorization**: Different access levels based on user roles
- **Password Security**: Secure password storage using bcrypt hashing

## API Endpoints

### Signup

**POST** `/api/auth/signup`

Request body:
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword",
  "phone_number": "1234567890",
  "role": "customer"
}
```

### Login

**POST** `/api/auth/login`

Request body:
```json
{
  "email": "john@example.com",
  "password": "securepassword"
}
```

### Get User Profile (Protected Route)

**GET** `/api/auth/profile`

Headers:
```
Authorization: Bearer <jwt_token>
```

## How to Run

1. Make sure PostgreSQL is installed and running
2. Update database configuration in `cmd/main.go` as needed
3. Run the service:

```bash
go run cmd/main.go
```

The service will start on port 8080 by default.

## Database Schema

The service uses a PostgreSQL database with the following schema:

```sql
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    mfa_secret VARCHAR(255),
    phone_number VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT false,
    phone_verified BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sessions (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## Integration with Other Services

To integrate with this authentication service from other services:

1. Send login/signup requests to the appropriate endpoints
2. Store the JWT token received from successful authentication
3. Include the token in the Authorization header for subsequent requests
4. Use the token claims to verify user role and permissions

## Project Structure

- `cmd/`: Main application entry point
- `pkg/models/`: Data models and request/response structures
- `pkg/database/`: Database connection and repositories
- `pkg/auth/`: Authentication service and HTTP handlers
- `pkg/utils/`: Utilities for password hashing, token generation, etc.