CREATE TABLE `user` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',

    -- 用户信息
    `uid` int(11) NOT NULL COMMENT '用户唯一ID',
    `nickname` varchar(255) NOT NULL COMMENT '昵称',
    `avatar` varchar(500) DEFAULT NULL COMMENT '头像',

    -- 公共参数
    `country` varchar(255) NULL COMMENT '国家',
    `ip` varchar(255) NULL COMMENT 'IP',
    `user_agent` varchar(512) NULL COMMENT 'UserAgent',

    -- 扩展信息
    `last_login_at` DATETIME(3) DEFAULT NULL COMMENT '最后登录时间',
    `login_count` int DEFAULT '0' COMMENT '登录次数',
    
    `status` int(11) NOT NULL DEFAULT 1 COMMENT '状态', -- 1: 正常, 2: 禁用
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';