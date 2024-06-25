CREATE TABLE userAuthentication (
    email VARCHAR(256),
    password_hash VARCHAR(60),
    last_login timestamp,
    PRIMARY KEY(email)
);


