-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS Contact (
  id        INT(11)         NOT NULL AUTO_INCREMENT,
  name      VARCHAR(191)    NOT NULL,
  ctime     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY (name),
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE Contact;
