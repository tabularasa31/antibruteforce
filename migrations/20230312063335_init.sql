-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS lists
(
    subnet      cidr primary key,
    list_type   text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS lists;
-- +goose StatementEnd
