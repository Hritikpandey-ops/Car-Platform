#   Car Platform – Microservices Architecture

This project is a microservices-based backend system built with Go, Docker, and PostgreSQL. It supports modular, scalable components for a car-related platform such as user authentication, vehicle management, document handling, and user profiling.

---

##   Services Implemented

### 1. **Auth Service (`auth-service`)**
Handles user registration, login, and JWT authentication.

-    Endpoints:
  - `POST /register` – Register a new user.
  - `POST /login` – Authenticate a user and return a JWT.
  - `GET /api/me` – Token verification endpoint.
-    Role-Based Access Control (RBAC): Supports `user`, `admin`.

---

### 2. **Vehicle Service (`vehicle-service`)**
CRUD operations for vehicle data, secured by token verification through `auth-service`.

-    Endpoints:
  - `POST /vehicles` – Add a vehicle.
  - `GET /vehicles` – List all vehicles.
  - `GET /vehicles/:id` – Get vehicle by ID.
  - `PATCH /vehicles/:id` – Update vehicle.
  - `DELETE /vehicles/:id` – Delete vehicle.

-   Middleware:
  - Token verification by calling `auth-service`'s `/api/me`.

---

### 3. **Document Service (`document-service`)**
Handles uploading and associating documents (e.g., RCs, images) with vehicles. Stores files in **MinIO** and metadata in **PostgreSQL**.

-    Endpoints:
  - `POST /documents/upload` – Upload a file.
  - `GET /documents/vehicle/:id` – Get docs for a vehicle.
  - `PATCH /documents/:id` – Update doc metadata.
  - `DELETE /documents/:id` – Delete a document.

-  MinIO is used as object storage.
-  JWT token verification via `auth-service`.

---

### 4. **User Service (`user-service`)**
Handles user profile management. Auth data (email/password) is stored in `auth-service`, while additional profile details are managed here.

-  Endpoints:
  - `POST /users` – Create a user profile.
  - `GET /users` – Get all profiles.
  - `GET /users/:id` – Get a specific user profile.
  - `PATCH /users/:id` – Update profile.
  - `DELETE /users/:id` – Delete profile.

---

## Tech Stack

| Tech            | Description                            |
|-----------------|----------------------------------------|
| Go              | Programming language for all services  |
| Gin             | Web framework                          |
| PostgreSQL      | Relational DB                          |
| MinIO           | Object storage                         |
| Docker          | Containerization                       |
| Docker Compose  | Multi-service orchestration            |
| JWT             | Token-based authentication             |

---

##  Database Design

### `auth-service` – `auth_db.users`
```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  is_verified BOOLEAN DEFAULT FALSE,
  verification_token TEXT,
  role VARCHAR(20) DEFAULT 'user'
);

-- vehicle service – vehicle_db.vehicles
sql
CREATE TABLE vehicles (
  id SERIAL PRIMARY KEY,
  brand TEXT NOT NULL,
  model TEXT NOT NULL,
  year INT NOT NULL,
  color TEXT NOT NULL,
  registration_number TEXT NOT NULL
);


-- document-service – vehicle_db.documents
sql
CREATE TABLE documents (
  id SERIAL PRIMARY KEY,
  filename TEXT NOT NULL,
  url TEXT NOT NULL,
  content_type TEXT NOT NULL,
  vehicle_id BIGINT,
  uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- user-service – user_db.users
sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  full_name TEXT,
  phone TEXT,
  address TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



##  Docker Setup
Run all services: bash
docker-compose up --build


Services & Ports:
Service	Port
Auth Service	   8081
Vehicle Service	8082
Document Service	8083
User Service	   8084
PostgreSQL	      5432
MinIO	9000 (API), 9001 (Console)

## Authentication
All protected services (vehicle, document, user) verify tokens using: http

GET http://auth-service:8081/api/me
Authorization: Bearer <token>

-- Testing
Use Postman or cURL to test endpoints. Include the JWT token in the Authorization header when required.

-- To Do (Next Steps)
 Implement Mail Service for sending verification emails.

 Add Admin Service for moderation & analytics.

 Integrate API Gateway (e.g., NGINX or Kong).

 Add Swagger / OpenAPI docs.

 Write unit & integration tests.

-- Author
Hritik Lalji Pandey