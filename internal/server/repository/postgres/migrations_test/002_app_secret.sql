-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS app_secret (
    id serial 		PRIMARY KEY,
    user_id  		INT,
    name 			VARCHAR (255) NOT NULL,
    data 			TEXT,
    description 	TEXT
);
--- add demo data
INSERT INTO app_secret (user_id, name, data, description) VALUES (1, 'name_1', 'data_1', 'description_1');
INSERT INTO app_secret (user_id, name, data, description) VALUES (1, 'name_2', 'data_2', 'description_2');
INSERT INTO app_secret (user_id, name, data, description) VALUES (2, 'name_3', 'data_3', 'description_3');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS app_secret;
-- +goose StatementEnd
