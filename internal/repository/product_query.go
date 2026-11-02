package repository

const saveProductQuery = `
INSERT INTO products (id, name, price_cents) VALUES (?, ?, ?)
`

const listProductsQuery = `
SELECT id, name, price_cents FROM products
`