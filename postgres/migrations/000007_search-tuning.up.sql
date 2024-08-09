CREATE INDEX character_name_search_idx ON player.character USING GIN (to_tsvector('simple', name));

CREATE INDEX corporation_name_search_idx ON player.corporation USING GIN (to_tsvector('simple', name));

CREATE INDEX alliance_name_search_idx ON player.alliance USING GIN (to_tsvector('simple', name));
