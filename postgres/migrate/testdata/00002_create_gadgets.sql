-- +goose Up
CREATE TABLE gadgets (id text PRIMARY KEY);

-- +goose Down
DROP TABLE gadgets;
