CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION fn_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
	IF OLD.* IS DISTINCT FROM NEW.* THEN
				NEW.updated_at = NOW(); 
	END IF;
	RETURN NEW;
END;
$$ language plpgsql ;
