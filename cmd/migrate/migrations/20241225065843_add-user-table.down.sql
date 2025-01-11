-- Drop triggers for automatically updating the 'updated_at' field on the users table
DROP TRIGGER IF EXISTS update_user_updated_at ON users;

-- Drop trigger functions
DROP FUNCTION IF EXISTS update_updated_at_users;

-- Drop indexes
DROP INDEX IF EXISTS idx_users_email;

-- Drop the users table (this will remove the dependency on the 'user_status' type)
DROP TABLE IF EXISTS users CASCADE;

-- Drop the user_status type (only after the users table is dropped)
DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_status') THEN
        DROP TYPE user_status;
    END IF;
END $$;
