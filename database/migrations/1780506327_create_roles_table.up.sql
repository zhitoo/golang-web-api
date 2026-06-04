CREATE TABLE roles (
    id SERIAL PRIMARY KEY,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP WITH TIME ZONE NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    created_by INTEGER NULL,
    modified_by INTEGER NULL,
    deleted_by INTEGER NULL,

    name VARCHAR(15) NOT NULL UNIQUE,

    CONSTRAINT fk_roles_created_by
        FOREIGN KEY (created_by) REFERENCES users(id),

    CONSTRAINT fk_roles_modified_by
        FOREIGN KEY (modified_by) REFERENCES users(id),

    CONSTRAINT fk_roles_deleted_by
        FOREIGN KEY (deleted_by) REFERENCES users(id)
);