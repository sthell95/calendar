CREATE TABLE IF NOT EXISTS events (
    id uuid PRIMARY KEY UNIQUE,
    title varchar(100),
    timezone varchar(40) default 'Europe/Riga',
    date_from timestamp,
    date_to timestamp,
    user_id uuid REFERENCES users(id)
);