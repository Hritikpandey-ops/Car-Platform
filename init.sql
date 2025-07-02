-- Create all required databases
CREATE DATABASE auth_db;
CREATE DATABASE vehicle_db;
CREATE DATABASE document;
CREATE DATABASE user_db;

-- ========================
-- Schema for auth_db
-- ========================
\connect auth_db

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  is_verified BOOLEAN DEFAULT FALSE,
  verification_token TEXT
);

ALTER TABLE users ADD COLUMN role VARCHAR(20) DEFAULT 'user';

-- ========================
-- Schema for vehicle_db
-- ========================
\connect vehicle_db

CREATE TABLE vehicles (
  id SERIAL PRIMARY KEY,
  brand TEXT NOT NULL,
  model TEXT NOT NULL,
  year INT NOT NULL,
  color TEXT NOT NULL,
  registration_number TEXT NOT NULL
);

-- ========================
-- Schema for document DB
-- ========================
\connect document

CREATE TABLE documents (
  id SERIAL PRIMARY KEY,
  filename TEXT NOT NULL,
  url TEXT NOT NULL,
  content_type TEXT NOT NULL,
  vehicle_id INT,
  uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

--=========================
-- Schema Key Constraints
--=========================
\connect user_db

CREATE TABLE user_profiles (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL, -- maps to `users.id` in auth_db
    full_name TEXT,
    phone TEXT,
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

