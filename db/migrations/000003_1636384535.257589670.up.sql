ALTER TABLE events RENAME COLUMN date_from TO time;
ALTER TABLE events RENAME COLUMN date_to TO duration;
ALTER TABLE events ALTER COLUMN duration SET DATA TYPE integer USING 0;
ALTER TABLE events ALTER COLUMN notes type varchar [];

CREATE TABLE IF NOT EXISTS notes (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    note varchar(60) NOT NULL,
    event_id uuid REFERENCES events(id) ON DELETE CASCADE
);