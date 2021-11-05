CREATE TABLE IF NOT EXISTS events (
    id uuid PRIMARY KEY UNIQUE,
    title varchar(100),
    description text,
    timezone varchar(40) default 'Europe/Riga',
    date_from timestamp,
    date_to timestamp,
    notes varchar ARRAY,
    user_id uuid REFERENCES users(id)
);