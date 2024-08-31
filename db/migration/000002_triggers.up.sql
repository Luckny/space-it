CREATE OR REPLACE FUNCTION set_updated_columns()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER on_update_set_updated_columns
  BEFORE UPDATE
  ON permissions
  FOR EACH ROW
  EXECUTE PROCEDURE set_updated_columns();
