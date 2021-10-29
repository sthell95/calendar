CREATE TABLE IF NOT EXISTS users
(
    id       uuid PRIMARY KEY,
    login    varchar(40) UNIQUE,
    password varchar(40),
    salt varchar(10) UNIQUE,
    timezone varchar(30)
)