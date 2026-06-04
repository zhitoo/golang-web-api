CREATE TABLE users (
    id SERIAL PRIMARY KEY,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP WITH TIME ZONE NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    created_by INTEGER NULL,
    modified_by INTEGER NULL,
    deleted_by INTEGER NULL,

    username VARCHAR(20) NOT NULL UNIQUE,
    first_name VARCHAR(30) NULL,
    last_name VARCHAR(30) NULL,

    mobile_number VARCHAR(11) UNIQUE NULL,
    email VARCHAR(60) UNIQUE NULL,

    password VARCHAR(128) NOT NULL,

    enabled BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT fk_users_created_by
        FOREIGN KEY (created_by) REFERENCES users(id),

    CONSTRAINT fk_users_modified_by
        FOREIGN KEY (modified_by) REFERENCES users(id),

    CONSTRAINT fk_users_deleted_by
        FOREIGN KEY (deleted_by) REFERENCES users(id)
);