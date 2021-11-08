ALTER TABLE events RENAME COLUMN date_from TO time;
ALTER TABLE events RENAME COLUMN date_to TO duration;
ALTER TABLE events ALTER COLUMN duration TYPE time;