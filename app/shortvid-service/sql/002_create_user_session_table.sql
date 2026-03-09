/* 创建用户会话表 */
CREATE TABLE user_session (
    `id` bigint(32) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',
    `user_uid` int(11) NOT NULL COMMENT '用户id',
    `session_id` varchar(255) NOT NULL COMMENT '会话id',
    `ip` varchar(255) NOT NULL COMMENT 'ip',
    `user_agent` varchar(255) NOT NULL COMMENT '用户代理',
    `expires_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '过期时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_session_id` (`session_id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户会话表';