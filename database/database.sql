CREATE DATABASE IF NOT EXISTS `hbl`;

USE `hbl`;

CREATE TABLE IF NOT EXISTS `addresses` (
  `address` VARBINARY(16) NOT NULL,
  `comment` VARCHAR(100) NOT NULL,
  `expires_at` TIMESTAMP NOT NULL,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE INDEX `idx_address` (`address`),
  PRIMARY KEY (`address`)
);

-- Procedures

DELIMITER //

-- CALL addresses()
-- Lists IP addresses.
DROP PROCEDURE IF EXISTS addresses //
CREATE PROCEDURE addresses()
BEGIN
  SELECT
    INET_NTOA(address) AS address, comment, expires_at, created_at, updated_at
  FROM
    addresses
  ORDER BY
    created_at DESC;
END //
