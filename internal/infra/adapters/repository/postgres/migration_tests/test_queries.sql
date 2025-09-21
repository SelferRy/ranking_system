SELECT schema_name FROM information_schema.schemata WHERE schema_name = 'ranking_system';
SELECT table_name FROM information_schema.tables WHERE table_schema = 'ranking_system';
SELECT * FROM ranking_system.banners WHERE id = 1;

INSERT INTO ranking_system.banners (id, description, banned_at, created_at, updated_at)
VALUES (1, 'test banner', NULL, NOW(), NULL);

INSERT INTO ranking_system.banners (id, description)
VALUES (2, 'test banner');

INSERT INTO ranking_system.slots (id, description) VALUES (1, 'test slot') RETURNING *;
SELECT * FROM ranking_system.slots;

SELECT * FROM ranking_system.banner_slot;