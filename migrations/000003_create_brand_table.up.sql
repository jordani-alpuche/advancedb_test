CREATE TABLE brand(  
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    brand_name varchar(255) NOT NULL,
    brand_description text,
    brand_created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    brand_updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);