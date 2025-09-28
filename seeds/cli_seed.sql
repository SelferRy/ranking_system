-- seeds/cli_seed.sql
SET search_path = public, ranking_system;

WITH b1 AS (
    INSERT INTO ranking_system.banners(description) VALUES ('Seed Banner A') RETURNING id
), b2 AS (
    INSERT INTO ranking_system.banners(description) VALUES ('Seed Banner B') RETURNING id
), s AS (
    INSERT INTO ranking_system.slots(description) VALUES ('Seed Slot CLI') RETURNING id
), g AS (
    INSERT INTO ranking_system.groups(description) VALUES ('Seed Group CLI') RETURNING id
), ins1 AS (
    INSERT INTO ranking_system.banner_slot (banner_id, slot_id)
        SELECT b1.id, s.id FROM b1, s
        UNION ALL
        SELECT b2.id, s.id FROM b2, s
)
INSERT INTO ranking_system.banner_stats (banner_id, slot_id, group_id, impressions, clicks)
SELECT b1.id, s.id, g.id, 100, 5 FROM b1, s, g
UNION ALL
SELECT b2.id, s.id, g.id, 200, 12 FROM b2, s, g;