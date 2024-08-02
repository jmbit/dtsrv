CREATE TABLE IF NOT EXISTS containersession (
  id INTEGER PRIMARY KEY,
  container TEXT,
  session TEXT
);

CREATE TABLE IF NOT EXISTS containertimeout (
  id INTEGER PRIMARY KEY,
  created INTEGER --Creation timestamp in UNIX time
  last_accessed INTEGER --Last access timestamp in UNIX time
)
