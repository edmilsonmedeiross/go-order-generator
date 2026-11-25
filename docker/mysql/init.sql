CREATE DATABASE IF NOT EXISTS products_db;
GRANT ALL PRIVILEGES ON products_db.* TO 'orders_user'@'%';
USE orders_db;
SOURCE /migrations/001_create_orders.sql;
SOURCE /migrations/004_create_outbox_events.sql;
USE products_db;
SOURCE /migrations/002_create_products.sql;
