-- Trigger function shared by all resource tables
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS offers (
    uuid       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    username   TEXT NOT NULL,
    type       VARCHAR(100) NOT NULL,
    name       VARCHAR(100) NOT NULL,
    price      NUMERIC(10,2) NOT NULL CHECK (price >= 0 AND price <= 1000)
);

CREATE TRIGGER offers_updated_at
    BEFORE UPDATE ON offers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
