-- Trigger function shared by all resource tables
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS boardgames (
    id            SERIAL PRIMARY KEY,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name          VARCHAR(100) NOT NULL UNIQUE,
    publisher     VARCHAR(100) NOT NULL,
    player_number INTEGER NOT NULL CHECK (player_number BETWEEN 1 AND 16),
    boardgame_id  INTEGER REFERENCES boardgames(id) ON DELETE SET NULL
);

CREATE TRIGGER boardgames_updated_at
    BEFORE UPDATE ON boardgames
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS tags (
    name       VARCHAR(30) PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER tags_updated_at
    BEFORE UPDATE ON tags
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS categories (
    name       VARCHAR(30) PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS mechanisms (
    name       VARCHAR(30) PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER mechanisms_updated_at
    BEFORE UPDATE ON mechanisms
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS boardgame_tags (
    boardgame_id INTEGER NOT NULL REFERENCES boardgames(id) ON DELETE CASCADE,
    tag_name     VARCHAR(30) NOT NULL REFERENCES tags(name) ON DELETE CASCADE,
    PRIMARY KEY (boardgame_id, tag_name)
);

CREATE TABLE IF NOT EXISTS boardgame_categories (
    boardgame_id  INTEGER NOT NULL REFERENCES boardgames(id) ON DELETE CASCADE,
    category_name VARCHAR(30) NOT NULL REFERENCES categories(name) ON DELETE CASCADE,
    PRIMARY KEY (boardgame_id, category_name)
);

CREATE TABLE IF NOT EXISTS boardgame_mechanisms (
    boardgame_id   INTEGER NOT NULL REFERENCES boardgames(id) ON DELETE CASCADE,
    mechanism_name VARCHAR(30) NOT NULL REFERENCES mechanisms(name) ON DELETE CASCADE,
    PRIMARY KEY (boardgame_id, mechanism_name)
);