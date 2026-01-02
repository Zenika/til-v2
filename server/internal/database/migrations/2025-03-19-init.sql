--
-- Init database with default tables. Drop existing tables if any.
-- Created on 2025-03-19 by Florian Forestier, Jim Gloaguen
--

PRAGMA
foreign_keys = ON;

DROP TABLE IF EXISTS `post_tags`;
DROP TABLE IF EXISTS `user_favorites_posts`;
DROP TABLE IF EXISTS `posts`;
DROP TABLE IF EXISTS `users`;

CREATE TABLE `users`
(
    `id`             TEXT    NOT NULL PRIMARY KEY,
    `display_name`   TEXT    NOT NULL,
    `google_id`      TEXT    NOT NULL UNIQUE,
    `default_filter` TEXT    NOT NULL        DEFAULT '',
    `is_admin`       BOOLEAN NOT NULL        DEFAULT FALSE,
    `feed_key`       TEXT    NOT NULL UNIQUE DEFAULT ''
);

CREATE TABLE `posts`
(
    `id`            TEXT NOT NULL PRIMARY KEY,
    `creation_date` DATETIME      DEFAULT CURRENT_TIMESTAMP NOT NULL,
    `user_id`       TEXT NOT NULL,
    `title`         TEXT NOT NULL,
    `link`          TEXT NOT NULL DEFAULT '',
    `content`       TEXT NOT NULL DEFAULT '',
    `image`         TEXT          DEFAULT '',
    `icon`          TEXT          DEFAULT '',
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE `post_tags`
(
    `post_id` TEXT NOT NULL,
    `tag`     TEXT NOT NULL,
    PRIMARY KEY (post_id, tag),
    FOREIGN KEY (post_id) REFERENCES posts (id)
);

CREATE TABLE `user_favorites_posts`
(
    `user_id` TEXT NOT NULL,
    `post_id` TEXT NOT NULL,
    PRIMARY KEY (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (post_id) REFERENCES posts (id)
);
