-- turns out these numbers get pretty big
ALTER TABLE player.corporation
ALTER COLUMN shares TYPE bigint USING shares::bigint;
