# API Documentation

This directory contains API documentation in various formats.

## Bruno Collection

The `bruno` folder contains the complete API specification in [Bruno](https://www.usebruno.com/) format - a modern API testing and documentation tool.

### Quick Start

1. Open Bruno application
2. Click "Open Collection"
3. Navigate to `bruno/User Management API`
4. Start testing API endpoints

### What's Included

- **7 Request Endpoints** - All CRUD operations for users plus health checks
- **Environment Configuration** - Pre-configured base URL and API version
- **Inline Documentation** - Each request includes detailed docs with examples
- **Error Handling** - Complete error code reference

### Directory Structure

```
bruno/
└── User Management API/
    ├── bruno.json              # Collection metadata
    ├── environment.bru         # Environment variables
    ├── README.md               # Complete documentation
    ├── Users/
    │   ├── Create User.bru
    │   ├── List All Users.bru
    │   ├── Get User by ID.bru
    │   ├── Update User.bru
    │   └── Delete User.bru
    └── System/
        ├── Health Check.bru
        └── Ready Check.bru
```

## API Overview

**Service**: User Management API  
**Framework**: Go + Gin  
**Base URL**: `http://localhost:8080/api/v1`

### Core Endpoints

| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/users` | Create new user |
| GET | `/users` | List all users |
| GET | `/users/:id` | Get specific user |
| PUT | `/users/:id` | Update user |
| DELETE | `/users/:id` | Delete user |

### System Endpoints

| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/health` | Health check |
| GET | `/ready` | Readiness check |

## Using the Documentation

### In Bruno

1. Import the collection
2. Customize environment variables (base URL, credentials, etc.)
3. Run requests individually or create test scenarios
4. View response bodies and status codes
5. Check inline documentation for each endpoint

### Exporting

Bruno collections can be exported for sharing with team members or importing into other tools.

## Next Steps

- See `bruno/User Management API/README.md` for detailed endpoint documentation
- Test the health check first to verify the API is running
- Follow the suggested testing workflow for a complete CRUD example

