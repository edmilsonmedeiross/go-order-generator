package repository

const saveOrderQuery = `
INSERT INTO orders (id, status, customer_id, name) VALUES (?, ?, ?, ?)
`
const saveOrderItemQuery = `
INSERT INTO order_items (id, order_id, product_id, quantity, price_cents, product_name) VALUES (?, ?, ?, ?, ?, ?)
`
const getProductInfoQuery = `
SELECT id, name, price_cents FROM products WHERE id = ?
`
const saveOutboxEventQuery = `
INSERT INTO outbox_events (id, event_type, aggregate_id, payload) VALUES (?, ?, ?, ?)
`