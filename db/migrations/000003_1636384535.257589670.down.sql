ALTER TABLE events RENAME COLUMN time TO date_from;
ALTER TABLE events RENAME COLUMN duration TO date_to;
-- ALTER TABLE events ALTER COLUMN date_to TYPE timestamp with time zone using date_to at time zone 'UTC';
ALTER TABLE events ADD COLUMN notes varchar ARRAY;
DROP TABLE notes
