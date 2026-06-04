CREATE TABLE role_user (
    role_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,

    PRIMARY KEY (role_id, user_id),

    CONSTRAINT fk_role_user_role_id
        FOREIGN KEY (role_id) REFERENCES roles(id),

    CONSTRAINT fk_role_user_user_id
        FOREIGN KEY (user_id) REFERENCES users(id)
);