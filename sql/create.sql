CREATE DATABASE IF NOT EXISTS feed_system
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_general_ci;

USE feed_system;

CREATE TABLE IF NOT EXISTS users (
  id BIGINT NOT NULL COMMENT 'snowflake user id',
  username VARCHAR(64) NOT NULL COMMENT 'login account',
  password_hash VARCHAR(100) NOT NULL COMMENT 'bcrypt password hash',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id),
  UNIQUE KEY uk_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='users';
