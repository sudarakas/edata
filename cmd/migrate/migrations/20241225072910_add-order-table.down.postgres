-- Drop triggers for automatically updating the 'updated_at' field on the orders table
DROP TRIGGER IF EXISTS update_order_updated_at ON orders;

-- Drop trigger functions
DROP FUNCTION IF EXISTS update_updated_at_orders;

-- Drop indexes
DROP INDEX IF EXISTS idx_subscriptions_name;
DROP INDEX IF EXISTS idx_orders_user_id;

-- Drop tables that depend on the 'order_status' type first
DROP TABLE IF EXISTS order_items CASCADE;
DROP TABLE IF EXISTS orders CASCADE;

-- Drop the order_status type (only after dependent tables are dropped)
DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
        DROP TYPE order_status;
    END IF;
END $$;
