-- PostgreSQL uchun maqolalar bazasi va jadvallarni yaratish skripti (Clean Architecture)

-- UUID ishlatish uchun kengaytmani yoqish (agar oldin yoqilmagan bo'lsa)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Maqolalar jadvali (Articles) - Barcha ma'lumotni jamlagan yagona jadval
CREATE TABLE IF NOT EXISTS articles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    access_type VARCHAR(50) NOT NULL, -- ochiq yoki yopiq
    abstract TEXT,
    authors JSONB NOT NULL DEFAULT '[]'::jsonb, -- Mualliflar massivi JSONB formatida bitta ustunda
    journal VARCHAR(255),
    publisher VARCHAR(255),
    publisher_date VARCHAR(50), 
    doi VARCHAR(100) UNIQUE,
    url VARCHAR(255),
    pdf_url VARCHAR(255),
    source_url VARCHAR(255),
    key_words TEXT[] -- PostgreSQL uchun massiv (Kalit so'zlar uchun)
);
