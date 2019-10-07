-- GENERATE TABLE FOR REMOTE CONFIG --
-- COMMENT DROP TABLE TO DISABLE DEBUGGING
DROP TABLE users, tribes, keys, tribe_assign, key_shares;

-- User data
CREATE TABLE users (
	user_id SERIAL PRIMARY KEY,
	first_name VARCHAR(20) NOT NULL,
	last_name VARCHAR(20) NOT NULL,
	email VARCHAR(50) NOT NULL,
	password VARCHAR(255) NOT NULL,
	role INT NOT NULL DEFAULT 0,
	UNIQUE(email)
);

-- Tribe data
CREATE TABLE tribes (
	tribe_id SERIAL PRIMARY KEY,
	tribe_name VARCHAR(50) NOT NULL,
	lead_id INT NOT NULL REFERENCES users(user_id),
	description VARCHAR(200),
	UNIQUE(tribe_name)
);

-- Key data
CREATE TABLE keys (
	key_id SERIAL PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	key_name VARCHAR (50) NOT NULL,
	key_value VARCHAR (300) NOT NULL,
	key_type VARCHAR (15) NOT NULL DEFAULT "STRING",
	description VARCHAR (200),
	platform VARCHAR (50) NOT NULL,
	expire_date TIMESTAMP NOT NULL DEFAULT NOW() + (5 * interval '1 week'),
	user_id INT NOT NULL REFERENCES users(user_id),
	tribe_id INT NOT NULL REFERENCES tribes(tribe_id),
	app_version VARCHAR (20) NOT NULL,
	status VARCHAR (20) NOT NULL
);

-- Table for assigning members
CREATE TABLE tribe_assign (
	tribe_id INT NOT NULL REFERENCES tribes(tribe_id),
	user_id INT NOT NULL REFERENCES users(user_id),
	PRIMARY KEY (tribe_id, user_id)
);

-- Table that stores shared key pointers
CREATE TABLE key_shares (
	key_id INT NOT NULL REFERENCES keys(key_id),
	user_id INT NOT NULL REFERENCES users(user_id),
	PRIMARY KEY (key_id, user_id)
);


-- key edit log is too expensive
-- use json data to store

-- -- Table that logs key edits
-- CREATE TABLE key_edits (
-- 	keyedit_id SERIAL PRIMARY KEY,
-- 	key_id INT NOT NULL REFERENCES keys(key_id),
-- 	user_id INT NOT NULL REFERENCES users(user_id),
-- 	edit_timestamp TIMESTAMP NOT NULL,
-- 	prev_value VARCHAR(200) NOT NULL,
-- 	curr_value VARCHAR(200) NOT NULL
-- );
