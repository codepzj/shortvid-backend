CREATE TABLE account(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',

    `user_uid` int(11) NOT NULL COMMENT '用户唯一ID',

    `email` varchar(64) NULL COMMENT '邮箱',
    `password` varchar(32) NULL COMMENT '密码',

    -- 第三方登录信息
    `provider` varchar(50) DEFAULT 'firebase' COMMENT '主要登录提供商',
    `provider_uid` varchar(100) DEFAULT NULL COMMENT '第三方平台用户UID',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_uid_provider` (`user_uid`, `provider`),
    UNIQUE KEY `uk_provider_provider_id` (`provider`, `provider_uid`)
)