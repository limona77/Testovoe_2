
CREATE TABLE users (
                       id serial PRIMARY KEY UNIQUE,
                       email varchar(255) NOT NULL UNIQUE,
                       password varchar(100) NOT NULL,
                       birthday date NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE tokens
(
    id serial PRIMARY KEY UNIQUE ,
    refresh_token varchar NOT NULL ,
    user_id  serial UNIQUE,
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES "users" (id)
);

CREATE TABLE subscriptions (
    user_id INT NOT NULL,
    subscribed_to_id INT NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id),

);

CREATE INDEX IF NOT EXISTS idx_subscriptions_user ON subscriptions (user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_subscribed_to ON subscriptions (subscribed_to_id);
