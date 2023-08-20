DROP TABLE IF EXISTS news;

CREATE TABLE news (
    id BIGSERIAL PRIMARY KEY, -- первичный ключ
    title TEXT NOT NULL,
	content TEXT,
	pubtime BIGINT NOT NULL,
	link TEXT
);
