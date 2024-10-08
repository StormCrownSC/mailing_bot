CREATE SCHEMA IF NOT EXISTS admin;

CREATE TABLE IF NOT EXISTS admin.admins (
    telegram_id BIGINT UNIQUE PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE NOT NULL
);

CREATE SCHEMA IF NOT EXISTS mailing;

CREATE TABLE IF NOT EXISTS mailing.users (
    telegram_id BIGINT UNIQUE PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE NOT NULL
);

CREATE TABLE IF NOT EXISTS mailing.texts (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE NOT NULL
);

-- for future updates
-- CREATE TABLE IF NOT EXISTS mailing.date (
--
-- );