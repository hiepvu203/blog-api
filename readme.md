# BLOG-API

A RESTful API for a blogging platform built with Golang and Gin. This project provides endpoints for user authentication, post management, category management, and comment features. It supports JWT authentication, role-based access control, and pagination, making it suitable for building modern blog or CMS applications.

---

## Features

- User registration, login, profile, password change, and self-deletion
- Role-based access: admin and client
- Admin management for users, posts, and categories
- CRUD operations for posts, categories, and comments
- JWT authentication middleware
- Pagination for listing resources
- Error handling with descriptive messages

---

## Tech Stack

- **Language:** Golang
- **Framework:** Gin Gonic
- **Database:** PostgreSQL, configure in `.env`
- **ORM:** GORM
- **Authentication:** JWT

---

## Project Structure

```
.
├── cmd/
│   └── main.go                # Application entry point
├── internal/
│   ├── config/                # Database and app configuration
│   ├── controllers/           # HTTP handlers for each resource
│   ├── dto/                   # Data transfer objects (request/response)
│   ├── entities/              # Database models
│   ├── repositories/          # Data access layer
│   ├── routes/                # Route definitions
│   └── services/              # Business logic
├── pkg/
│   ├── middlewares/           # JWT, admin, and other middlewares
│   └── utils/                 # Utility functions
├── doc/
│   ├── blog-api-doc.md        # API documentation (Markdown)
│   └── blog-api-doc.json      # Postman collection
├── .env                       # Environment variables
├── go.mod / go.sum            # Go modules
└── README.md                  # Project documentation
```

---

## Getting Started

### Prerequisites

- Go 1.18+
- A running database (PostgreSQL)

### Installation

1. **Clone the repository:**
    ```sh
    git clone https://github.com/hiepvu203/blog-api.git
    cd blog-api
    ```

2. **Configure environment variables:**
    - Copy `.env.example` to `.env` and update database credentials and JWT secret.

3. **Install dependencies:**
    ```sh
    go mod tidy
    ```

4. **Run database migrations** (if any).

5. **Start the server:**
    ```sh
    go run cmd/main.go
    ```
    The API will be available at `http://localhost:9090`.

---

## API Documentation

- **Postman Collection:** Import [`doc/blog-api-doc.json`](doc/blog-api-doc.json) into Postman

### Example Endpoints

- `POST /users/register` – Register a new user
- `POST /users/login` – User login, returns JWT token
- `GET /users/me` – Get current user profile (requires JWT)
- `GET /posts` – List all posts (with pagination)
- `POST /posts` – Create a new post (requires authentication)
- `PUT /posts/:id` – Update a post (owner or admin)
- `DELETE /posts/:id` – Delete a post (owner or admin)
- `GET /categories` – List all categories
- `POST /posts/:postId/comments` – Add a comment to a post
- `PUT /comments/:id` – Update a comment (owner or admin)
- `GET /posts/:postId/comments` – Get comments for a post

---

## Authentication

- Use JWT tokens in the `Authorization` header:  
  `Authorization: Bearer <token>`
- Admin-only endpoints require a user with `role: admin`.

---
