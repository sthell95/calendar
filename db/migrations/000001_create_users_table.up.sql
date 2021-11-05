CREATE TABLE IF NOT EXISTS users
(
    id       uuid PRIMARY KEY,
    login    varchar(40) UNIQUE,
    password varchar(60),
    timezone varchar(30)
)