CREATE TABLE Users (
    guid UUID PRIMARY KEY,
    email VARCHAR(254) UNIQUE,
    ip INET,
    refreshtoken BYTEA
);

CREATE INDEX ip_index ON Users USING BTREE (ip);