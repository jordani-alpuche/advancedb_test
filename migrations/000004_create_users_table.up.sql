CREATE TABLE users(  
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    first_name varchar(50) NOT NULL,
    last_name varchar(50) NOT NULL,
    username varchar(50) NOT NULL,
    password varchar(255) NOT NULL,
    email varchar(100) NOT NULL,
    phone_number varchar(15),
    role varchar(20) NOT NULL    
);