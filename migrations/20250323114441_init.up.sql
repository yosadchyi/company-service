CREATE TABLE IF NOT EXISTS companies
(
    id                  UUID PRIMARY KEY,
    name                VARCHAR(15) UNIQUE NOT NULL,
    description         TEXT,
    amount_of_employees INT                NOT NULL,
    registered          BOOLEAN            NOT NULL,
    type                VARCHAR(50)        NOT NULL
);
