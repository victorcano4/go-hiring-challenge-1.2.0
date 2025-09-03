CREATE TABLE IF NOT EXISTS product_categories (
    id SERIAL PRIMARY KEY,
    code VARCHAR(32) UNIQUE NOT NULL,
    name VARCHAR(255)
);

INSERT INTO product_categories (code, name) VALUES
('CLOTHING', '123'),
('SHOES', '456'),
('ACCESSORIES', '789');

ALTER TABLE products
ADD COLUMN category_id INT,
ADD CONSTRAINT fk_category
FOREIGN KEY (category_id) REFERENCES product_categories (id);

UPDATE products
SET category_id = (SELECT id FROM product_categories WHERE code = 'CLOTHING')
WHERE code IN ('PROD001', 'PROD004', 'PROD007');

UPDATE products
SET category_id = (SELECT id FROM product_categories WHERE code = 'SHOES')
WHERE code IN ('PROD002', 'PROD006');

UPDATE products
SET category_id = (SELECT id FROM product_categories WHERE code = 'ACCESSORIES')
WHERE code IN ('PROD003', 'PROD005', 'PROD008');