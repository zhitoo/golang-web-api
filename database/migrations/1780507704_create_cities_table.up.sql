CREATE TABLE cities (
    id SERIAL PRIMARY KEY,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP WITH TIME ZONE NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    created_by INTEGER NULL,
    modified_by INTEGER NULL,
    deleted_by INTEGER NULL,

    country_id INTEGER NOT NULL,
    name VARCHAR(30) NOT NULL,
   

    CONSTRAINT fk_cities_country_id
        FOREIGN KEY (country_id) REFERENCES countries(id),


    CONSTRAINT fk_cities_created_by
        FOREIGN KEY (created_by) REFERENCES users(id),

    CONSTRAINT fk_cities_modified_by
        FOREIGN KEY (modified_by) REFERENCES users(id),

    CONSTRAINT fk_cities_deleted_by
        FOREIGN KEY (deleted_by) REFERENCES users(id)
);