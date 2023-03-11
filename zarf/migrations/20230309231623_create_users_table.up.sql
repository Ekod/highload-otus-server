CREATE TABLE IF NOT EXISTS users
(
    id         INTEGER PRIMARY KEY AUTO_INCREMENT,
    first_name VARCHAR(30) NOT NULL,
    last_name  VARCHAR(50) NOT NULL,
    age        INTEGER(120) NOT NULL,
    gender     ENUM ('male', 'female') NOT NULL,
    interests  VARCHAR(200) NOT NULL,
    city       VARCHAR(70) NOT NULL,
    email      VARCHAR(50) NOT NULL,
    password   VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
