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
    ```
    PORT=...
    DB_HOST=...
    DB_PORT=...
    DB_USER=...
    DB_PASSWORD=...
    DB_NAME=...
    JWT_SECRET=...
    ```

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

[![Swagger Hub](https://img.shields.io/badge/docs-SwaggerHub-blue)](https://app.swaggerhub.com/apis/vxunhip/blog-api/1.0)

- **SwaggerHub:** Xem và thử nghiệm API trực tiếp trên [SwaggerHub](https://app.swaggerhub.com/apis/vxunhip/blog-api/1.0).