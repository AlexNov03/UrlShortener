-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS url (
    url_id SERIAL PRIMARY KEY, 
    short_url CHAR(10) UNIQUE NOT NULL, 
    original_url VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS url;
-- +goose StatementEnd
