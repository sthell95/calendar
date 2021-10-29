CREATE TABLE IF NOT EXISTS users
(
    id       uuid PRIMARY KEY,
    login    varchar(40) UNIQUE,
    timezone varchar(30)
)