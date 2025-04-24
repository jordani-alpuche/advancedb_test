CREATE TABLE products(  
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    product_name VARCHAR(255) NOT NULL,
    product_description TEXT,
    product_price FLOAT NOT NULL,
    product_category_id INT NOT NULL,
    product_brand_id INT NOT NULL,
    product_qty INT NOT NULL,
    product_status VARCHAR(50) NOT NULL,
    product_create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    product_update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    product_purchased_from VARCHAR(255)
);