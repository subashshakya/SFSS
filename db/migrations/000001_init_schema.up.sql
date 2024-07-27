CREATE TABLE SecureFile (
    id TEXT PRIMARY KEY,
    filename TEXT NOT NULL,
    file_data BYTEA NOT NULL,
    original_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),,
    user_id INT NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES User(id)
);


CREATE TABLE SuperSecret (
    id TEXT PRIMARY KEY,
    secret TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    user_id INT NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES User(id)
);


CREATE TABLE User (
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    middle_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    phone_number TEXT NOT NULL
);


CREATE TABLE FileSharing (
    id SERIAL PRIMARY KEY,
    file_id TEXT NOT NULL,
    sender_id INT NOT NULL,
    recipient_id INT NOT NULL,
    shared_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_file
        FOREIGN KEY(file_id)
        REFERENCES SecureFile(id),
    CONSTRAINT fk_sender
        FOREIGN KEY(sender_id)
        REFERENCES User(id),
    CONSTRAINT fk_recipient
        FOREIGN KEY(recipient_id)
        REFERENCES User(id)
);


CREATE TABLE SecretSharing (
    id SERIAL PRIMARY KEY,
    secret_id TEXT NOT NULL,
    sender_id INT NOT NULL,
    recipient_id INT NOT NULL,
    shared_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_secret
        FOREIGN KEY(secret_id)
        REFERENCES SuperSecret(id)
    CONSTRAINT fk_sender
        FOREIGN KEY(sender_id)
        REFERENCES User(id),
    CONSTRAINT fk_receiver
        FOREIGN KEY(recipient_id)
        REFERENCES User(id)
);


CREATE TABLE SecretFileCount (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES User(id)
);

CREATE TABLE SecretPasswordCount (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES User(id)
);
