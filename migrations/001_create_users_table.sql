-- migrations/001_create_users_table.sql
-- 用户表，存储账号、密码（bcrypt hash）和角色
CREATE TABLE IF NOT EXISTS `users` (
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3)     DEFAULT NULL,
    `updated_at` datetime(3)     DEFAULT NULL,
    `deleted_at` datetime(3)     DEFAULT NULL,
    `username`   varchar(64)     NOT NULL,
    `password`   varchar(256)    NOT NULL,
    `role`       varchar(32)     DEFAULT 'user',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_users_username` (`username`),
    KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- casbin_rule 表由 gorm-adapter 自动创建，结构参考如下：
-- CREATE TABLE IF NOT EXISTS `casbin_rule` (
--     `id`    bigint unsigned NOT NULL AUTO_INCREMENT,
--     `ptype` varchar(100)    DEFAULT NULL,
--     `v0`    varchar(100)    DEFAULT NULL,
--     `v1`    varchar(100)    DEFAULT NULL,
--     `v2`    varchar(100)    DEFAULT NULL,
--     `v3`    varchar(100)    DEFAULT NULL,
--     `v4`    varchar(100)    DEFAULT NULL,
--     `v5`    varchar(100)    DEFAULT NULL,
--     PRIMARY KEY (`id`),
--     UNIQUE KEY `idx_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
