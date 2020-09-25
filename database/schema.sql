-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CREATE TABLE "base_table" (created_at TIMESTAMP NOT NULL, updated_at TIMESTAMP NOT NULL)

CREATE TABLE "user_account" (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
) INHERITS (base_table);

CREATE TABLE collection (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
    title VARCHAR(255) NOT NULL,
    number_of_chapters INTEGER,
    most_recent_upload_date TIMESTAMP,
    user_id uuid,
    FOREIGN KEY (user_id) REFERENCES "user_account"(id)
) INHERITS (base_table);

CREATE TABLE chapter (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    number_of_pages INTEGER,
    chapter_number INTEGER,
    collection_id uuid,
    FOREIGN KEY (collection_id) REFERENCES "collection"(id)
) INHERITS (base_table);

CREATE TABLE page (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
    page BYTEA,
    title VARCHAR(255) NOT NULL,
    chapter_id uuid,
    order_number INTEGER NOT NULL,
    FOREIGN KEY (chapter_id) REFERENCES "chapter"(id)
) INHERITS (base_table);