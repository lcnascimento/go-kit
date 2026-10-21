-- +goose Up
CREATE TABLE widgets (id text PRIMARY KEY);

-- +goose Down
DROP TABLE widgets;
