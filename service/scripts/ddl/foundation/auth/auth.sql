CREATE TABLE `user`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `username`   varchar(256)        NOT NULL COMMENT '用户名',
    `nick_name`  varchar(256)        NOT NULL COMMENT '用户昵称',
    `auth_type`  varchar(128)        NOT NULL COMMENT '认证类型',
    `password`   varchar(256) DEFAULT NULL COMMENT '用户名',
    `email`      varchar(128) DEFAULT NULL COMMENT '邮箱',
    `phone`      varchar(64)  DEFAULT NULL COMMENT '手机号',
    `avatar`     varchar(256) DEFAULT NULL COMMENT '头像路径',
    `status`     tinyint(4)   DEFAULT 1 COMMENT '用户状态，1:未激活, 2:已激活',
    `created_at` datetime            NOT NULL COMMENT '创建时间',
    `updated_at` datetime            NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `deleted_at` datetime     DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`) COMMENT '用户名索引',
    KEY `uk_nick_name` (`nick_name`) COMMENT '用户昵称索引',
    UNIQUE KEY `uk_email` (`email`) USING BTREE COMMENT '邮箱唯一索引',
    UNIQUE KEY `uk_phone` (`phone`) USING BTREE COMMENT '手机号唯一索引'
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE `namespace`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `namespace`   varchar(256)        NOT NULL COMMENT 'namespace',
    `description` varchar(1024)       NULL COMMENT '命名空间描述',
    `creator`     varchar(128)        NOT NULL COMMENT '创建人',
    `created_at`  datetime            NOT NULL COMMENT '创建时间',
    `updated_at`  datetime            NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `deleted_at`  datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_namespace` (`namespace`) COMMENT 'namespace索引'
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE `namespace_token`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name`       varchar(128)        NULL COMMENT '名称',
    `namespace`  varchar(256)        NOT NULL COMMENT 'namespace',
    `token`      varchar(256)        NOT NULL COMMENT 'token校验',
    `expired_at` datetime DEFAULT NULL COMMENT '失效时间',
    `creator`    varchar(128)        NOT NULL COMMENT '创建人',
    `created_at` datetime            NOT NULL COMMENT '创建时间',
    `updated_at` datetime            NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_namespace_token` (`namespace`, `token`) COMMENT '联合索引'
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8mb4;