CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS articles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    access_type VARCHAR(50) NOT NULL,
    abstract TEXT,
    authors JSONB NOT NULL DEFAULT '[]'::jsonb,
    journal VARCHAR(255),
    publisher VARCHAR(255),
    publisher_date VARCHAR(50), 
    doi VARCHAR(100) UNIQUE,
    url VARCHAR(255),
    pdf_url VARCHAR(255),
    source_url VARCHAR(255),
    key_words TEXT[]
);

CREATE TABLE IF NOT EXISTS article_views (
    article_id UUID PRIMARY KEY REFERENCES articles(id) ON DELETE CASCADE,
    views_count INT NOT NULL DEFAULT 0
);
