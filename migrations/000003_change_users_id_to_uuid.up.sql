-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Drop attendances table (will be recreated with UUID foreign key)
DROP TABLE IF EXISTS attendances;

-- Create a temporary table to store user data
CREATE TABLE users_temp AS SELECT * FROM users;

-- Drop the original users table
DROP TABLE users;

-- Recreate users table with UUID
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'employee',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Drop temporary table (data migration not possible from SERIAL to UUID)
DROP TABLE users_temp;
