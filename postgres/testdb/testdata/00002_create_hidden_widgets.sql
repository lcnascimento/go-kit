-- +goose Up
CREATE TABLE hidden_widgets (id text PRIMARY KEY);
INSERT INTO hidden_widgets (id) VALUES ('w-1');
ALTER TABLE hidden_widgets ENABLE ROW LEVEL SECURITY;
ALTER TABLE hidden_widgets FORCE ROW LEVEL SECURITY;
CREATE POLICY hide_every_row ON hidden_widgets USING (false);

-- +goose Down
DROP TABLE hidden_widgets;
