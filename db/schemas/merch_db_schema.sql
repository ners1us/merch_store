CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255)        NOT NULL,
    coins    INTEGER             NOT NULL DEFAULT 1000
);

CREATE TABLE IF NOT EXISTS merches
(
    name  VARCHAR(255) PRIMARY KEY,
    price INTEGER NOT NULL
);

INSERT INTO merches (name, price)
VALUES ('t-shirt', 80),
       ('cup', 20),
       ('book', 50),
       ('pen', 10),
       ('powerbank', 200),
       ('hoody', 300),
       ('umbrella', 200),
       ('socks', 10),
       ('wallet', 50),
       ('pink-hoody', 500);

CREATE TABLE IF NOT EXISTS coin_transfers
(
    id           SERIAL PRIMARY KEY,
    from_user_id INTEGER   NOT NULL REFERENCES users (id),
    to_user_id   INTEGER   NOT NULL REFERENCES users (id),
    amount       INTEGER   NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS purchases
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER      NOT NULL REFERENCES users (id),
    merch_item VARCHAR(255) NOT NULL REFERENCES merches (name),
    created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);
