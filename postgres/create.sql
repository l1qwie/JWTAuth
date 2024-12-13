CREATE TABLE Users (
    guid VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    ip VARCHAR(255),
    refreshtoken BYTEA
);
