CREATE TABLE `users` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',

    -- 用户信息
    `user_uid` varchar(255) NOT NULL COMMENT '用户唯一ID',
    `nickname` varchar(255) NOT NULL COMMENT '昵称',
    `avatar` varchar(500) DEFAULT NULL COMMENT '头像',
    `email` varchar(255) DEFAULT NULL COMMENT '邮箱',

    -- 第三方登录信息
    `provider` varchar(50) DEFAULT 'firebase' COMMENT '主要登录提供商',
    `provider_uid` varchar(100) DEFAULT NULL COMMENT '第三方平台用户UID',
  
    -- 扩展信息
    `last_login_at` datetime(3) DEFAULT NULL COMMENT '最后登录时间',
    `login_count` int DEFAULT '0' COMMENT '登录次数',
    
    `status` int(11) NOT NULL DEFAULT 1 COMMENT '状态',
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_uid` (`user_uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';