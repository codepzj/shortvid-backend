-- 创建视频表
CREATE TABLE video (
    `id` bigint(32) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    -- 视频底层
    `vgroup` varchar(255) NOT NULL COMMENT '视频底层分组',

    -- 视频信息
    `user_uid` int(11) NOT NULL COMMENT '用户id',
    `description` varchar(255) NOT NULL COMMENT '视频描述',
    `category` varchar(255) NOT NULL COMMENT '视频分类',
    `tags` varchar(255) NOT NULL COMMENT '视频标签',
    `custom_tags` varchar(255) NOT NULL COMMENT '自定义标签',

    -- 视频统计
    `like_count` int(11) NOT NULL DEFAULT 0 COMMENT '点赞数',
    `view_count` int(11) NOT NULL DEFAULT 0 COMMENT '浏览数',
    
    -- 视频状态
    `status` int(11) NOT NULL DEFAULT 1 COMMENT '状态 1-上传中, 2-处理中, 3-审核中, 4-发布成功, 5-拒绝发布, 6-系统失败, 7-封禁',

    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频表';