-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS app_user (
    id serial 			PRIMARY KEY,
    login 				VARCHAR (255) NOT NULL,
    password 			VARCHAR (255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS app_user;
-- +goose StatementEnd
