-- Create order_items table (Order details for the subscription order)
CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,                -- Order this item belongs to
    quantity INT NOT NULL DEFAULT 1,      -- Quantity of the subscription
    price DECIMAL(10, 2) NOT NULL,        -- Price of the subscription plan at the time of the order
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

-- Create index for order_id to optimize lookups
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items (order_id);

-- Add trigger for automatically updating the 'updated_at' field on any record change
CREATE OR REPLACE FUNCTION update_updated_at_order_items() 
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_order_item_updated_at
BEFORE UPDATE ON order_items
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_order_items();
