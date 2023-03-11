CREATE TABLE IF NOT EXISTS friends
(
    id        INTEGER PRIMARY KEY AUTO_INCREMENT,
    user_id   INTEGER NOT NULL,
    friend_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (friend_id) REFERENCES users (id)
);