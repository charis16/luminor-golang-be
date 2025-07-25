CREATE TABLE websites (
    id SERIAL PRIMARY KEY,
     uuid UUID UNIQUE DEFAULT gen_random_uuid(),
    about_us_brief_home_en TEXT,
    about_us_en TEXT,
    about_us_id TEXT,
    about_us_brief_home_id TEXT,
    address TEXT,
    phone_number VARCHAR(15),
    email VARCHAR(50),
    url_instagram VARCHAR(50),
    url_facebook VARCHAR(50),
    url_tiktok VARCHAR(50),
    video_web VARCHAR(255),
    video_mobile VARCHAR(255),
    meta_title VARCHAR(255),
    meta_desc TEXT,
    meta_keyword TEXT,
    og_image VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);