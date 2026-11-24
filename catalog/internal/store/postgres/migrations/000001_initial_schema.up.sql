-- Catalog initial schema

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE boardgames (
    id             SERIAL PRIMARY KEY,
    slug           VARCHAR(120) NOT NULL UNIQUE,
    name           VARCHAR(120) NOT NULL,
    description    TEXT NOT NULL DEFAULT '',
    year_published SMALLINT NOT NULL DEFAULT 0,
    min_players    SMALLINT NOT NULL,
    max_players    SMALLINT NOT NULL,
    min_play_time  SMALLINT NOT NULL DEFAULT 0,
    max_play_time  SMALLINT NOT NULL DEFAULT 0,
    min_age        SMALLINT,
    bgg_id         INTEGER UNIQUE,
    boardgame_id   INTEGER REFERENCES boardgames(id) ON DELETE SET NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ,
    CONSTRAINT boardgames_slug_not_empty CHECK (slug <> ''),
    CONSTRAINT boardgames_min_players_check CHECK (min_players BETWEEN 1 AND 16),
    CONSTRAINT boardgames_max_players_check CHECK (max_players BETWEEN 1 AND 16),
    CONSTRAINT boardgames_players_order_check CHECK (min_players <= max_players),
    CONSTRAINT boardgames_year_published_check CHECK (year_published = 0 OR year_published BETWEEN 1900 AND 2100),
    CONSTRAINT boardgames_min_play_time_check CHECK (min_play_time BETWEEN 0 AND 9999),
    CONSTRAINT boardgames_max_play_time_check CHECK (max_play_time BETWEEN 0 AND 9999),
    CONSTRAINT boardgames_play_time_order_check CHECK (min_play_time = 0 OR max_play_time = 0 OR min_play_time <= max_play_time)
);

CREATE TRIGGER boardgames_updated_at
    BEFORE UPDATE ON boardgames
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE categories (
    id         SERIAL PRIMARY KEY,
    slug       VARCHAR(100) NOT NULL UNIQUE,
    name       VARCHAR(100) NOT NULL,
    bgg_id     INTEGER UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT categories_slug_not_empty CHECK (slug <> '')
);

CREATE TRIGGER categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE mechanisms (
    id         SERIAL PRIMARY KEY,
    slug       VARCHAR(100) NOT NULL UNIQUE,
    name       VARCHAR(100) NOT NULL,
    bgg_id     INTEGER UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT mechanisms_slug_not_empty CHECK (slug <> '')
);

CREATE TRIGGER mechanisms_updated_at
    BEFORE UPDATE ON mechanisms
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE contributors (
    id         SERIAL PRIMARY KEY,
    slug       VARCHAR(100) NOT NULL UNIQUE,
    name       VARCHAR(100) NOT NULL,
    bio        TEXT,
    bgg_id     INTEGER UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT contributors_slug_not_empty CHECK (slug <> '')
);

CREATE TRIGGER contributors_updated_at
    BEFORE UPDATE ON contributors
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE boardgame_categories (
    boardgame_id INTEGER NOT NULL REFERENCES boardgames(id) ON DELETE CASCADE,
    category_id  INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (boardgame_id, category_id)
);

CREATE TABLE boardgame_mechanisms (
    boardgame_id INTEGER NOT NULL REFERENCES boardgames(id) ON DELETE CASCADE,
    mechanism_id INTEGER NOT NULL REFERENCES mechanisms(id) ON DELETE CASCADE,
    PRIMARY KEY (boardgame_id, mechanism_id)
);

CREATE TABLE boardgame_contributions (
    boardgame_id   INTEGER NOT NULL REFERENCES boardgames(id) ON DELETE CASCADE,
    contributor_id INTEGER NOT NULL REFERENCES contributors(id) ON DELETE CASCADE,
    role           VARCHAR(30) NOT NULL,
    credit_order   SMALLINT,
    PRIMARY KEY (boardgame_id, contributor_id, role)
);
