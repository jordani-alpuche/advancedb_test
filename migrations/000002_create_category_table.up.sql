CREATE TABLE category(  
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    category_name varchar(255) NOT NULL,
    category_description text,
    category_code varchar(50) NOT NULL,
    category_created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    category_updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);