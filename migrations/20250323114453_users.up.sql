CREATE TABLE IF NOT EXISTS users
(
    id       UUID PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT        NOT NULL
);

INSERT INTO users (id, username, password)
VALUES ('902B061B-9FE9-482D-A230-E760B03BBE39', 'admin', '$2a$10$3k5e47Xmp.duRXDlyW.8a.ZKHSY3q/LYHPeDQAli14/GYr1tvqxvm')
ON CONFLICT DO NOTHING;
