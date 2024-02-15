CREATE TABLE login_tokens (
  login_token uuid PRIMARY KEY,
  user_id uuid
) WITH default_time_to_live = 3600;
