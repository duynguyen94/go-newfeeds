CREATE TABLE `user` (
    id INT AUTO_INCREMENT PRIMARY KEY,
    hashed_password VARBINARY(255) NOT NULL,
    salt VARBINARY(255) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    dob DATE NOT NULL,
    email VARCHAR(100) NOT NULL,
    user_name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE `post` (
    id INT AUTO_INCREMENT PRIMARY KEY,
    fk_user_id INT,
    content_text TEXT NOT NULL,
    content_image_path VARCHAR(255) NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    visible BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (fk_user_id) REFERENCES `user`(id)
);

CREATE TABLE `comment` (
    id INT AUTO_INCREMENT PRIMARY KEY,
    fk_post_id INT,
    fk_user_id INT,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (fk_post_id) REFERENCES `post`(id),
    FOREIGN KEY (fk_user_id) REFERENCES `user`(id)
);

CREATE TABLE `like` (
    fk_post_id INT,
    fk_user_id INT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (fk_post_id, fk_user_id),
    FOREIGN KEY (fk_post_id) REFERENCES `post`(id),
    FOREIGN KEY (fk_user_id) REFERENCES `user`(id)
);

CREATE TABLE `user_user` (
    fk_user_id INT,
    fk_follower_id INT,
    PRIMARY KEY (fk_user_id, fk_follower_id),
    FOREIGN KEY (fk_user_id) REFERENCES `user`(id),
    FOREIGN KEY (fk_follower_id) REFERENCES `user`(id)
);