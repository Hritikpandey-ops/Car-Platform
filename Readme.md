# Auth Service – Microservice for Authentication

This is the **Auth Service** component of a microservices-based vehicle platform built using **Go (Golang)** and **Docker**. It handles user authentication and account management functionality such as signup, login, email verification, and JWT-based protected routes.

---

## Features

- **User Signup** with hashed password storage (bcrypt)
- **Login** with JWT token generation
- **Email Verification** via secure token
- **JWT Middleware** to protect internal routes
- **PostgreSQL** as the database
- Containerized with **Docker & Docker Compose**

---

## Tech Stack

- **Language**: Go (Golang)
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **Token Handling**: JWT
- **Email Delivery**: SMTP (Gmail App Password)
- **Containerization**: Docker & Docker Compose

---

## Project Structure
auth-service/
├── Dockerfile
├── .env
├── main.go
├── handlers/
│ └── auth.go
├── models/
│ └── user.go
├── database/
│ └── connect.go
├── utils/
│ ├── jwt.go
│ └── email.go
├── init.sql

---

## Getting Started

1. **Clone the repo**
   ```bash
   git clone https://github.com/Hritikpandey-ops/car-platform.git
   cd car-platform

Create a .env file in auth-service/:
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=auth_db

JWT_SECRET=your_jwt_secret

EMAIL_FROM=your_email@gmail.com
EMAIL_PASSWORD=your_app_password
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587


Start services with Docker
## docker-compose up --build


Test API Endpoints using Postman
## POST /signup
## GET /verify?token=...
## POST /login
## GET /api/me (with JWT)
