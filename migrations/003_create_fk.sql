ALTER TABLE order_items
ADD FOREIGN KEY (product_id) REFERENCES products(id);