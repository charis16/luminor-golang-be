CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE DEFAULT gen_random_uuid(),

    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    photo VARCHAR(255),
    description TEXT,

    password VARCHAR(100),      
    role VARCHAR(50),

    phone_number VARCHAR(15),

    url_instagram VARCHAR(100),
    url_tiktok VARCHAR(100),
    url_facebook VARCHAR(100),
    url_youtube VARCHAR(100),
    
    is_published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);