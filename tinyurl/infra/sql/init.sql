CREATE TABLE urls (
    shorten_url VARCHAR(256),
    original_url VARCHAR(2048),
    counter integer,
    expiration_date timestamp,
    PRIMARY KEY(shorten_url)
);
