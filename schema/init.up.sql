CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE,
    username TEXT UNIQUE,
    password_hash TEXT
);

CREATE TABLE IF NOT EXISTS session (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  token VARCHAR(255) NOT NULL,
  expiration_time TIMESTAMP NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS post (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  author_id INTEGER NOT NULL,
  like INTEGER DEFAULT 0,
  dislike INTEGER DEFAULT 0,
  title TEXT,
  category TEXT,
  content TEXT, 
  author TEXT,
  date TEXT, 
  FOREIGN KEY (author_id) REFERENCES user (id)
);

CREATE TABLE IF NOT EXISTS category (
  post_id INTEGER NOT NULL, 
  category TEXT,
  FOREIGN KEY (post_id) REFERENCES post (id)
);

CREATE TABLE IF NOT EXISTS comment(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		post_id INTEGER,
		like INTEGER DEFAULT 0,
		dislike INTEGER DEFAULT 0,
		text TEXT,
		author	TEXT, 
		date TEXT,
		FOREIGN KEY (post_id) REFERENCES post (id)
	);


CREATE TABLE IF NOT EXISTS reaction_post (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER,
  post_id INTEGER,
  like_is INTEGER DEFAULT(-1),
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  FOREIGN KEY (post_id) REFERENCES post (id) ON DELETE CASCADE
);


  CREATE TABLE IF NOT EXISTS reaction_comment (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER,
  comment_id INTEGER,
  like_is INTEGER DEFAULT(-1),
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  FOREIGN KEY (comment_id) REFERENCES comment (id) ON DELETE CASCADE
);