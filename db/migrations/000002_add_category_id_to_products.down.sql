ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_categories;
ALTER TABLE products DROP COLUMN IF EXISTS category_id;
