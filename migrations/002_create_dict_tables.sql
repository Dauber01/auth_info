-- migrations/002_create_dict_tables.sql
-- 字典类型表：定义可配置信息的分类（如数据威胁等级、社交媒体类型）
CREATE TABLE IF NOT EXISTS `dict_types` (
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at`  datetime(3)     DEFAULT NULL,
    `updated_at`  datetime(3)     DEFAULT NULL,
    `deleted_at`  datetime(3)     DEFAULT NULL,
    `code`        varchar(64)     NOT NULL    COMMENT '类型编码，全局唯一',
    `name`        varchar(128)    NOT NULL    COMMENT '类型名称',
    `description` varchar(256)    DEFAULT ''  COMMENT '描述',
    `sort`        int             DEFAULT 0   COMMENT '排序，数值越小越靠前',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_dict_types_code` (`code`),
    KEY `idx_dict_types_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  COMMENT = '字典类型表';

-- 字典数据表：存储各分类下的具体配置项
CREATE TABLE IF NOT EXISTS `dict_items` (
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at`  datetime(3)     DEFAULT NULL,
    `updated_at`  datetime(3)     DEFAULT NULL,
    `deleted_at`  datetime(3)     DEFAULT NULL,
    `type_code`   varchar(64)     NOT NULL    COMMENT '关联字典类型编码',
    `item_key`    varchar(64)     NOT NULL    COMMENT '配置项键',
    `item_value`  varchar(256)    NOT NULL    COMMENT '配置项值',
    `description` varchar(256)    DEFAULT ''  COMMENT '描述',
    `sort`        int             DEFAULT 0   COMMENT '排序，数值越小越靠前',
    `status`      tinyint         DEFAULT 1   COMMENT '状态：1 启用，0 禁用',
    PRIMARY KEY (`id`),
    KEY `idx_dict_items_type_code` (`type_code`),
    KEY `idx_dict_items_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  COMMENT = '字典数据表';

-- 预置字典类型
INSERT INTO `dict_types` (`code`, `name`, `description`, `sort`, `created_at`, `updated_at`) VALUES
    ('data_threat_level', '数据威胁等级', '描述数据资产的威胁风险级别', 1, NOW(), NOW()),
    ('social_media_type', '社交媒体类型', '支持接入的社交媒体平台分类', 2, NOW(), NOW());

-- 预置数据威胁等级
INSERT INTO `dict_items` (`type_code`, `item_key`, `item_value`, `description`, `sort`, `status`, `created_at`, `updated_at`) VALUES
    ('data_threat_level', '1', '低危',  '风险影响较小，可接受范围内',     1, 1, NOW(), NOW()),
    ('data_threat_level', '2', '中危',  '风险影响适中，需要关注并跟踪',   2, 1, NOW(), NOW()),
    ('data_threat_level', '3', '高危',  '风险影响较大，需要及时处置',     3, 1, NOW(), NOW()),
    ('data_threat_level', '4', '严重',  '风险影响极大，需要立即响应',     4, 1, NOW(), NOW());

-- 预置社交媒体类型
INSERT INTO `dict_items` (`type_code`, `item_key`, `item_value`, `description`, `sort`, `status`, `created_at`, `updated_at`) VALUES
    ('social_media_type', 'weibo',    '微博',     '新浪微博平台',        1, 1, NOW(), NOW()),
    ('social_media_type', 'wechat',   '微信',     '微信公众号/朋友圈',   2, 1, NOW(), NOW()),
    ('social_media_type', 'douyin',   '抖音',     '抖音短视频平台',      3, 1, NOW(), NOW()),
    ('social_media_type', 'twitter',  'Twitter',  'Twitter/X 平台',      4, 1, NOW(), NOW()),
    ('social_media_type', 'facebook', 'Facebook', 'Facebook 平台',       5, 1, NOW(), NOW());
