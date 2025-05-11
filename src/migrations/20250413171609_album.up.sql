CREATE TABLE albums (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE DEFAULT gen_random_uuid(),
    slug VARCHAR(255),
    title VARCHAR(255),
    category_id INT REFERENCES categories(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    description TEXT,
    images TEXT[], -- Postgres supports array data types, so this is valid
    thumbnail VARCHAR(255),
    is_published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);