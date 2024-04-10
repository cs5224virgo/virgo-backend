CREATE TABLE IF NOT EXISTS users (
  	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP,

	username VARCHAR UNIQUE NOT NULL,
	password VARCHAR NOT NULL,
	display_name VARCHAR
);

CREATE INDEX ON users (deleted_at);
CREATE INDEX ON users (username);