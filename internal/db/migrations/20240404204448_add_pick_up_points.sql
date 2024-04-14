-- +goose Up
-- +goose StatementBegin
CREATE TABLE pick_up_points
(
    id      BIGSERIAL PRIMARY KEY NOT NULL,
    name    TEXT                  NOT NULL DEFAULT '',
    address TEXT                  NOT NULL DEFAULT '',
    contact TEXT                  NOT NULL DEFAULT ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table pick_up_points;
-- +goose StatementEnd
