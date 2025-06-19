# BLOG-API Documentation

This document provides a detailed overview of the **BLOG-API**, a RESTful API for a blogging platform built with Golang and Gin. It supports user authentication, post management, category management, and comment features, with JWT authentication, role-based access control, and pagination.

---

## Table of Contents

1. [Admin Endpoints](#admin-endpoints)
   - [Users](#users)
     - [Get List of Users](#get-list-of-users)
     - [Change User Role](#change-user-role)
     - [Delete User](#delete-user)
2. [Post Endpoints](#post-endpoints)
   - [Get Posts](#get-posts)
3. [Comment Endpoints](#comment-endpoints)
   - [Add Comment](#add-comment)
   - [Update Comment](#update-comment)
   - [Get Comments for a Post](#get-comments-for-a-post)

---

## Admin Endpoints

### Users

#### Get List of Users

**Description**: Retrieves a list of users from the admin section of the application, including their ID, username, email, and role.

**Request**:
- **Method**: GET
- **URL**: `http://localhost:9090/admin/users`
- **Headers**:
  - `Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTA0OTUxNzMsInJvbGUiOiJhZG1pbiIsInVzZXJfaWQiOjR9.4-8f2TiLAFY87kS6cFl1a8tg0hpV8ByM9uMcjl4aWLM`

**Response**:
- **Status**: 200 OK
- **Content-Type**: `application/json; charset=utf-8`
- **Body**:
  ```json
  {
      "data": [
          {
              "id": 8,
              "username": "hiep03",
              "email": "vuxuanhiep1709@gmail.com",
              "role": "client"
          },
          {
              "id": 6,
              "username": "testuser01",
              "email": "test@example.com",
              "role": "client"
          },
          {
              "id": 4,
              "username": "testuser",
              "email": "test123@example.com",
              "role": "admin"
          },
          {
              "id": 1,
              "username": "admin",
              "email": "admin@example.com",
              "role": "admin"
          },
          {
              "id": 2,
              "username": "johndoe",
              "email": "john@example.com",
              "role": "admin"
          }
      ],
      "message": "Users fetched successfully.",
      "page": 1,
      "page_size": 10,
      "success": true,
      "total": 5
  }
  ```

**Error Responses**:
1. **Unauthorized (401)**:
   - **Body**: `{"error": "Authorization header required"}`
2. **Unauthorized (401, Invalid Token)**:
   - **Body**: `{"error": "Invalid token"}`
3. **Forbidden (403, Not Admin)**:
   - **Body**: `{"error": "Admin access required"}`
4. **Bad Request (400, Invalid Page)**:
   - **Body**: `{"success": false, "error": "Invalid page parameter"}`
5. **Bad Request (400, Invalid Page Size)**:
   - **Body**: `{"success": false, "error": "Invalid page_size parameter"}`

**Tests**:
- Ensures response status is 200.
- Validates response time is less than 200ms.
- Verifies response schema includes `data`, `message`, `page`, `page_size`, `success`, and `total`.
- Checks `data` is an array of user objects with `id`, `username`, `email`, and `role`.
- Ensures `success` is true.
- Confirms `data` array length is greater than 0.
- Generates a visualizer table for the response data.

| ID | Username   | Email                    | Role   |
|----|------------|--------------------------|--------|
| 8  | hiep03     | vuxuanhiep1709@gmail.com | client |
| 6  | testuser01 | test@example.com         | client |
| 4  | testuser   | test123@example.com      | admin  |
| 1  | admin      | admin@example.com        | admin  |
| 2  | johndoe    | john@example.com         | admin  |

#### Change User Role

**Description**: Updates the role of a specific user identified by their unique ID.

**Request**:
- **Method**: PUT
- **URL**: `http://localhost:9090/admin/users/6/role`
- **Headers**:
  - `Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTA0MDc5MTEsInJvbGUiOiJhZG1pbiIsInVzZXJfaWQiOjR9.ZmASocYH_GSUMQLYK1qj05FTveTJFjGoBife1TujUSI`
- **Body**:
  ```json
  { "role": "admin" }
  ```

**Response**:
- **Status**: 200 OK
- **Content-Type**: `application/json; charset=utf-8`
- **Body**:
  ```json
  {
      "success": true,
      "data": {
          "message": "User role updated"
      }
  }
  ```

**Error Responses**:
1. **Unauthorized (401)**:
   - **Body**: `{"error": "Authorization header required"}`
2. **Unauthorized (401, Invalid Token)**:
   - **Body**: `{"error": "Invalid token"}`
3. **Forbidden (403, Not Admin)**:
   - **Body**: `{"error": "Admin access required"}`
4. **Bad Request (400, Invalid User ID)**:
   - **Body**: `{"success": false, "error": "Invalid user id"}`
5. **Bad Request (400, No Value)**:
   - **Body**: `{"success": false, "error": "Key: 'Role' Error:Field validation for 'Role' failed on the 'required' tag"}`
6. **Bad Request (400, Incorrect Value)**:
   - **Body**: `{"success": false, "error": "Key: 'Role' Error:Field validation for 'Role' failed on the 'oneof' tag"}`
7. **Bad Request (400, User Not Found)**:
   - **Body**: `{"success": false, "error": "user not found"}`

**Tests**:
- Ensures response time is less than 200ms.
- Verifies response contains `success` and `data` keys.
- Validates response schema includes `success` (boolean) and `data` with `message` (string).
- Confirms `message` is a string.

#### Delete User

**Description**: Deletes a user from the system by their unique ID.

**Request**:
- **Method**: DELETE
- **URL**: `http://localhost:9090/admin/users/2`
- **Headers**:
  - `Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTA0OTg3OTksInJvbGUiOiJhZG1pbiIsInVzZXJfaWQiOjR9.wszw4EDNQt96ZDukiNe-wYEiua7JqrO4hI17yV9aptg`

**Response**:
- **Status**: 200 OK
- **Content-Type**: `application/json; charset=utf-8`
- **Body**:
  ```json
  {
      "success": true,
      "data": {
          "message": "User deleted successfully"
      }
  }
  ```

**Error Responses**:
1. **Unauthorized (401)**:
   - **Body**: `{"error": "Authorization header required"}`
2. **Unauthorized (401, Invalid Token)**:
   - **Body**: `{"error": "Invalid token"}`
3. **Forbidden (403, Not Admin)**:
   - **Body**: `{"error": "Admin access required"}`
4. **Bad Request (400, User Not Found)**:
   - **Body**: `{"success": false, "error": "Invalid user id"}`

---

## Post Endpoints

#### Get Posts

**Description**: Retrieves a list of posts with optional query parameters for pagination and filtering.

**Request**:
- **Method**: GET
- **URL**: `http://localhost:9090/posts`
- **Query Parameters** (optional):
  - `page`: Page number (e.g., `2`)
  - `page_size`: Number of results per page (e.g., `5`)
  - `title`: Search by title (e.g., `golang`)
  - `author`: Search by author (e.g., `hiep03`)
  - `category`: Search by category (e.g., `kien-thuc-cong-nghe`)

**Response**:
- **Status**: 200 OK
- **Content-Type**: `application/json; charset=utf-8`
- **Body** (example with pagination):
  ```json
  {
      "data": [
          {
              "id": 13,
              "title": "Quản lý session trong web app",
              "slug": "quan-ly-session-trong-web-app",
              "content": "Các kỹ thuật quản lý session trong ứng dụng web.",
              "thumbnail": "https://picsum.photos/seed/10/400/200",
              "category_id": 2,
              "category": "Lập trình",
              "author_id": 4,
              "author": "testuser",
              "status": "published",
              "created_at": "2025-06-18 13:32:01",
              "updated_at": "2025-06-18 13:32:01"
          },
          ...
      ],
      "page": 2,
      "page_size": 5,
      "success": true,
      "total": 16
  }
  ```

**Error Responses**:
1. **Not Found (200, No Articles)**:
   - **Body**: `{"data": null, "message": "No matching articles found.", "page": 1, "page_size": 10, "success": true, "total": 0}`
2. **Internal Server Error (500)**:
   - **Body**: `{"success": false, "error": "Could not fetch posts"}`

---

## Comment Endpoints

#### Add Comment

**Description**: Adds a comment to a specific post identified by its ID.

**Request**:
- **Method**: POST
- **URL**: `http://localhost:9090/posts/11/comments`
- **Headers**:
  - `Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTA1NjMwODcsInJvbGUiOiJhZG1pbiIsInVzZXJfaWQiOjR9.K8G5ZJHPlwXsU_5oHKQOUqjgHbD1PyRegMA1PDJFfs8`
- **Body**:
  ```json
  {
      "content": "Bài viết rất hay, cảm ơn bạn đã chia sẻ!"
  }
  ```

**Response**:
- **Status**: 201 Created
- **Content-Type**: `application/json; charset=utf-8`
- **Body**:
  ```json
  {
      "success": true,
      "data": {
          "message": "Comment created successfully"
      }
  }
  ```

**Error Responses**:
1. **Bad Request (400, Invalid Post ID)**:
   - **Body**: `{"success": false, "error": "Invalid post id"}`
2. **Unauthorized (401)**:
   - **Body**: `{"success": false, "error": "Authorization header required"}`

#### Update Comment

**Description**: Updates an existing comment identified by its unique ID.

**Request**:
- **Method**: PUT
- **URL**: `http://localhost:9090/comments/9`
- **Headers**:
  - `Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTA1NjMwODcsInJvbGUiOiJhZG1pbiIsInVzZXJfaWQiOjR9.K8G5ZJHPlwXsU_5oHKQOUqjgHbD1PyRegMA1PDJFfs8`
- **Body**:
  ```json
  {
      "content": "Nội dung comment mới"
  }
  ```

**Response**:
- **Status**: 200 OK
- **Content-Type**: `application/json; charset=utf-8`
- **Body**:
  ```json
  {
      "success": true,
      "data": {
          "message": "Comment updated successfully"
      }
  }
  ```

**Error Responses**:
1. **Not Found (404, Endpoint Not Found)**:
   - **Body**: `{"success": false, "error": "Endpoint not found"}`

#### Get Comments for a Post

**Description**: Retrieves comments associated with a specific post identified by its ID.

**Request**:
- **Method**: GET
- **URL**: `http://localhost:9090/posts/2/comments`
- **Headers**:
  - `Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTA0MzIzMTQsInJvbGUiOiJhZG1pbiIsInVzZXJfaWQiOjR9.RckYgrZBDbQSrB7JyBUQ526Q9nVUWYX4J3UeNLEw2Fg`
- **Body** (not typically used for GET, but included in last call):
  ```json
  {
      "content": "Nội dung comment mới"
  }
  ```

**Response**:
- **Status**: 200 OK
- **Content-Type**: `application/json; charset=utf-8`
- **Body**:
  ```json
  {
      "success": true,
      "data": {
          "comments": [
              {
                  "id": 7,
                  "post_id": 2,
                  "user_id": 4,
                  "content": "Bài viết rất hay, cảm ơn bạn đã chia sẻ!",
                  "created_at": "2025-06-17T22:18:57.752769Z",
                  "updated_at": "2025-06-17T22:18:57.752769Z"
              }
          ],
          "page": 1,
          "page_size": 10,
          "total": 1
      }
  }
  ```

**Error Responses**:
1. **Bad Request (400, Invalid Post ID)**:
   - **Body**: `{"success": false, "error": "Invalid post id"}`

---

## Notes
- All admin endpoints require JWT authentication with admin role.
- Pagination is supported for user and post lists.
- Ensure valid JWT tokens are used to avoid unauthorized errors.
- The API uses UTF-8 encoding for responses, supporting multilingual content.