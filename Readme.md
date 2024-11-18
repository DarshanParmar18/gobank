# GoBanking API

A secure REST API for banking operations built with Go, featuring JWT authentication and PostgreSQL database integration.

## Features

- User account creation and management
- Secure login with JWT authentication
- Account balance operations
- Money transfer between accounts
- RESTful API endpoints
- PostgreSQL database storage
- Middleware for request authentication

## Tech Stack

- Go 1.23.1
- Gorilla Mux (HTTP router)
- JWT for authentication
- PostgreSQL
- Testify for testing
- Crypto for password hashing

## API Endpoints

### Public Endpoints

- `POST /login` - Authenticate user and get JWT token
- `POST /account` - Create new account

### Protected Endpoints (Requires JWT)

- `GET /account` - Get all accounts
- `GET /account/{id}` - Get account by ID
- `DELETE /account/{id}` - Delete account
- `POST /transfer` - Transfer money between accounts

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Protected endpoints require a valid JWT token in the `x-jwt-token` header.

## Installation

1. Clone the repository:

```bash
git clone https://github.com/darshanparmar18/gobank.git
```

2. Install dependencies:

```bash
go mod download
```

3. Set up environment variables:

```bash
export JWT_SECRET=your_jwt_secret_key
export DB_URL=your_postgres_connection_string
```

4. Run the server:

```bash
make run
```

## Development

To build the project:

```bash
make build
```

To run tests:

```bash
make test
```

## License

MIT License

## Contact

Darshan Parmar - [@darshanparmar18](https://github.com/darshanparmar18)
Project Link: [https://github.com/darshanparmar18/gobank](https://github.com/darshanparmar18/gobank)
