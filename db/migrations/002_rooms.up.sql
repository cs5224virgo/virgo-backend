CREATE TABLE IF NOT EXISTS rooms (
  	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP,

    code VARCHAR UNIQUE NOT NULL,
    name VARCHAR NOT NULL,
    description VARCHAR
);

CREATE INDEX ON rooms (deleted_at);

CREATE TABLE IF NOT EXISTS rooms_users_memberships (
    room_id INTEGER NOT NULL REFERENCES rooms (id),
    user_id INTEGER NOT NULL REFERENCES users (id),
    unread INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX ON rooms_users_memberships (room_id);
CREATE INDEX ON rooms_users_memberships (user_id);