-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS author
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_author_timestamp() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE OR REPLACE TRIGGER trigger_update_author_timestamp
    BEFORE UPDATE
    ON author
    FOR EACH ROW
EXECUTE FUNCTION update_author_timestamp();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_update_author_timestamp ON author;
DROP FUNCTION IF EXISTS update_author_timestamp;
DROP TABLE IF EXISTS author;

