-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS ranking_system;

CREATE TABLE IF NOT EXISTS ranking_system.banners (
    id BIGSERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    banned_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ranking_system.slots (
    id BIGSERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    banned_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ranking_system.groups (
    id BIGSERIAL PRIMARY KEY,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS ranking_system.banner_slot (
    banner_id BIGINT NOT NULL REFERENCES ranking_system.banners(id) ON DELETE CASCADE,
    slot_id BIGINT NOT NULL REFERENCES ranking_system.slots(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    PRIMARY KEY (banner_id, slot_id)
);

CREATE TABLE IF NOT EXISTS ranking_system.banner_stats (
    banner_id   BIGINT NOT NULL,
    slot_id     BIGINT NOT NULL,
    group_id    BIGINT NOT NULL,
    impressions BIGINT NOT NULL DEFAULT 0,
    clicks      BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (banner_id, slot_id, group_id),
    FOREIGN KEY (banner_id, slot_id) REFERENCES ranking_system.banner_slot (banner_id, slot_id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES ranking_system.groups(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ranking_system.banner_stats;
DROP TABLE IF EXISTS ranking_system.banner_slot;
DROP TABLE IF EXISTS ranking_system.groups;
DROP TABLE IF EXISTS ranking_system.slots;
DROP TABLE IF EXISTS ranking_system.banners;

DROP SCHEMA IF EXISTS ranking_system;
-- +goose StatementEnd