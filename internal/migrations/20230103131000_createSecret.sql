-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS app_secret (
    id serial 		PRIMARY KEY,
    user_id  		INT,
    name 			VARCHAR (255) NOT NULL,
    data 			TEXT,
    description 	TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE app_secret IF EXISTS;
-- +goose StatementEnd
