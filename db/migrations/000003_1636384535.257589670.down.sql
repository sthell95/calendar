ALTER TABLE events RENAME COLUMN time TO date_from;
ALTER TABLE events RENAME COLUMN duration TO date_to;
ALTER TABLE events ALTER COLUMN date_to TYPE timestamp;