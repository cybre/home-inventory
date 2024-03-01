CREATE TYPE user_households_room (
  household_id UUID,
  room_id UUID,
  name TEXT,
  item_count INT
);

CREATE TABLE user_households (
  user_id TEXT,
  household_id UUID,
  name TEXT,
  location TEXT,
  description TEXT,
  item_count INT,
  rooms LIST<FROZEN<user_households_room>>,
  timestamp TIMESTAMP,
  PRIMARY KEY (user_id, household_id)
);
