CREATE TYPE user_households_room (
  household_id UUID,
  room_id UUID,
  name TEXT,
  sort_order INT,
  tstamp TIMESTAMP
);

CREATE TABLE user_households (
  user_id TEXT,
  household_id UUID,
  name TEXT,
  location TEXT,
  description TEXT,
  rooms MAP<TEXT, frozen<user_households_room>>,
  tstamp TIMESTAMP,
  sort_order INT,
  PRIMARY KEY (user_id, household_id)
);