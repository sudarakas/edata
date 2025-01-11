-- Drop triggers for automatically updating the 'updated_at' field
DROP TRIGGER IF EXISTS update_subscription_updated_at ON subscriptions;

-- Drop trigger functions
DROP FUNCTION IF EXISTS update_updated_at_subscriptions;

-- Drop indexes
DROP INDEX IF EXISTS idx_subscriptions_name;

-- Drop tables
DROP TABLE IF EXISTS subscriptions CASCADE;
