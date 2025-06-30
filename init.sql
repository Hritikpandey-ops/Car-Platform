-- Create auth_db
CREATE DATABASE auth_db;

-- Create vehicle_db
CREATE DATABASE vehicle_db;

-- Switch to auth_db and create tables
\connect auth_db

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  is_verified BOOLEAN DEFAULT FALSE,
  verification_token TEXT
);

ALTER TABLE users ADD COLUMN role VARCHAR(20) DEFAULT 'user';


-- Switch to vehicle_db and create tables
\connect vehicle_db

CREATE TABLE vehicles (
  id SERIAL PRIMARY KEY,
  brand TEXT NOT NULL,
  model TEXT NOT NULL,
  year INT NOT NULL,
  color TEXT NOT NULL,
  registration_number TEXT NOT NULL
);
