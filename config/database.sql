CREATE DATABASE IF NOT EXISTS `hbl`;

USE `hbl`;

CREATE TABLE IF NOT EXISTS `addresses` (
  `ip` VARBINARY(16) NOT NULL,
  `author` VARCHAR(100) NOT NULL,
  `action` VARCHAR(100) NOT NULL,
  `comment` VARCHAR(100) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE INDEX `idx_ip` (`ip`),
  PRIMARY KEY (`ip`)
);

CREATE TABLE IF NOT EXISTS `abuseipdb_metadata` (
  `ip` VARBINARY(16),
  `country_code` VARCHAR(50),
  `usage_type` VARCHAR(50),
  `isp` VARCHAR(50),
  `abuse_confidence_score` INT,
  `num_distinct_users` INT,
  `total_reports` INT,
  `last_reported_at` TIMESTAMP,
  UNIQUE INDEX `idx_ip` (`ip`),
  PRIMARY KEY (`ip`)
);
