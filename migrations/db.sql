CREATE TABLE IF NOT EXISTS Banners
(
    banner_id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    feature_id INT NOT NULL,
    content JSON NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS Banner_tags
(
    banner_id INT NOT NULL REFERENCES Banners(banner_id) ON DELETE CASCADE,
    tag_id INT NOT NULL
);

CREATE INDEX IF NOT EXISTS feature_idx ON Banners(feature_id);
CREATE INDEX IF NOT EXISTS tag_idx ON Banner_tags(tag_id);