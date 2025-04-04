DO $$
BEGIN

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE IF NOT EXISTS users (
    user_id uuid UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name VARCHAR(150) NOT NULL,
    last_name  VARCHAR(150) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password VARCHAR(150) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS user_email ON users(email);

CREATE TABLE IF NOT EXISTS recording (
    id SERIAL PRIMARY KEY,
    user_id uuid UNIQUE NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    uploaded BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT recording_id UNIQUE (user_id, id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE INDEX IF NOT EXISTS recording_id ON recording(id);

END
$$