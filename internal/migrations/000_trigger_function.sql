-- +migrate Up
DO $trigger_function$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_proc WHERE proname = 'update_updated_at_column'
    ) THEN
        EXECUTE 'CREATE FUNCTION update_updated_at_column()
        RETURNS TRIGGER AS $$ 
        BEGIN
            NEW.updated_at = CURRENT_TIMESTAMP;
            RETURN NEW;
        END;
        $$ LANGUAGE plpgsql';
    END IF;
END
$trigger_function$;

-- +migrate Down
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE; 