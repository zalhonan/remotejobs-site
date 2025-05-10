-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS jobs_raw (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    title TEXT,
    content_pure TEXT,
    source_link VARCHAR(2048) NOT NULL,
    main_technology VARCHAR(255),
    stop_words TEXT[],
    date_posted TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    date_parsed TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_jobs_raw_main_technology ON jobs_raw(main_technology);

CREATE TABLE IF NOT EXISTS telegram_channels (
    id BIGSERIAL PRIMARY KEY,
    tag VARCHAR(255) NOT NULL UNIQUE,
    last_post_id BIGINT DEFAULT 0,
    date_channel_added TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    posts_parsed BIGINT NOT NULL DEFAULT 0,
    date_last_parsed TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_telegram_channels_tag ON telegram_channels(tag);

CREATE TABLE IF NOT EXISTS technologies (
    id BIGSERIAL PRIMARY KEY,
    technology VARCHAR(100) NOT NULL UNIQUE,
    keywords TEXT[] NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0 CHECK (sort_order BETWEEN -100 AND 100),
    count BIGINT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_technologies_technology ON technologies(technology);
CREATE INDEX IF NOT EXISTS idx_technologies_keywords ON technologies USING GIN (keywords);

CREATE TABLE IF NOT EXISTS stop_words (
    id BIGSERIAL PRIMARY KEY,
    word VARCHAR(255) NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_stop_words_word ON stop_words(word);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS technologies;
DROP TABLE IF EXISTS telegram_channels;
DROP TABLE IF EXISTS jobs_raw;
DROP TABLE IF EXISTS stop_words;
-- +goose StatementEnd
