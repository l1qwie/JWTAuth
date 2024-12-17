CREATE TABLE Users (
    guid UUID PRIMARY KEY,
    email VARCHAR(254) UNIQUE,
    ip INET,
    refreshtoken BYTEA
);