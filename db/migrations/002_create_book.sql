-- +goose Up
CREATE TABLE IF NOT EXISTS book
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_book_timestamp() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE OR REPLACE TRIGGER trigger_update_book_timestamp
    BEFORE UPDATE
    ON book
    FOR EACH ROW
EXECUTE FUNCTION update_book_timestamp();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_update_book_timestamp ON book;
DROP FUNCTION IF EXISTS update_book_timestamp;
DROP TABLE IF EXISTS book;

