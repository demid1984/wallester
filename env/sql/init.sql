CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    birthdate timestamp NOT NULL,
    gender VARCHAR(16) NOT NULL,
    email VARCHAR(255) NOT NULL,
    address VARCHAR(200) NOT NULL,
    version int NOT NULL DEFAULT 1
);

CREATE INDEX name_idx ON customers(last_name, first_name);

INSERT INTO customers(first_name, last_name, birthdate, gender, email, address) VALUES
    ('Петр', 'Петров', NOW() - INTERVAL '20 YEARS', 'Male', 'petroff@test.ee', 'Viru väljak 4-6, 10111 Tallinn'),
    ('Сидов', 'Сидоров', NOW() - INTERVAL '25 YEARS', 'Male', 'sidor2000@test.ee', 'Kaubamaja 8, 10145 Tallinn'),
    ('Мария', 'Сидорова', NOW() - INTERVAL '26 YEARS', 'Female', 'marymaria@test.ee', 'Liivalaia 33, 10118 Tallinn');