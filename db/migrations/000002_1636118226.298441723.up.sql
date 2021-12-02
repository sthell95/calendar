CREATE TABLE IF NOT EXISTS events (
    id uuid PRIMARY KEY UNIQUE default gen_random_uuid(),
    title varchar(100) NOT NULL,
    description text,
    timezone varchar(40) default 'Europe/Riga',
    date_from timestamp with time zone NOT NULL ,
    date_to timestamp with time zone NOT NULL,
    notes varchar ARRAY,
    "user" uuid REFERENCES users(id)
);