CREATE TYPE message_type AS ENUM (
  'normal',
  'system',
  'summary'
);
-- ALTER TYPE enum_type ADD VALUE 'new_value'; -- for future reference

CREATE TABLE IF NOT EXISTS messages (
  	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP,

    content VARCHAR NOT NULL,
    type message_type NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users (id),
    room_id INTEGER NOT NULL REFERENCES rooms (id)
);

CREATE INDEX ON messages (deleted_at);
CREATE INDEX ON messages (room_id);