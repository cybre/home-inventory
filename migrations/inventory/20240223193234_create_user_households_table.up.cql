CREATE TABLE user_households (
  user_id TEXT,
  household_id UUID,
  name TEXT,
  location TEXT,
  description TEXT,
  PRIMARY KEY (user_id, household_id)
);
