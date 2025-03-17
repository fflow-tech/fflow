CREATE TABLE `function`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name`          varchar(256) NOT NULL COMMENT '函数名',
    `namespace`     varchar(256) NOT NULL COMMENT '命名空间',
    `creator`       varchar(128) NOT NULL COMMENT '创建人',
    `updater`       varchar(128) NOT NULL COMMENT '更新人',
    `code`          text COMMENT '代码',
    `token`         varchar(256)  NOT NULL COMMENT 'token',
    `description`   varchar(1024) DEFAULT NULL COMMENT '描述',
    `language`      varchar(256) NOT NULL COMMENT '所使用语言 javascript/golang/starlark',
    `version`       int(8) NOT NULL COMMENT '版本号',
    `created_at`    datetime     NOT NULL COMMENT '创建时间',
    `updated_at`    datetime     NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `deleted_at`    datetime      DEFAULT NULL COMMENT '删除时间',
    `input_schema`  json          DEFAULT NULL COMMENT '函数入参格式',
    `output_schema` json          DEFAULT NULL COMMENT '函数返回值格式',
    PRIMARY KEY (`id`) USING BTREE COMMENT '主键索引',
    KEY             `idx_creator` (`creator`) COMMENT '创建者索引',
    KEY             `idx_namespace` (`namespace`) USING BTREE COMMENT '命名空间索引',
    KEY             `idx_name_version` (`name`,`version`) USING BTREE COMMENT '函数名版本索引',
    KEY             `idx_created_at` (`created_at`) COMMENT '创建时间索引',
    KEY             `idx_deleted_at` (`deleted_at`) COMMENT '删除时间索引'
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `run_history`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `namespace`  varchar(256) NOT NULL COMMENT '命名空间',
    `name`       varchar(256) NOT NULL COMMENT '函数名',
    `operator`   varchar(128) NOT NULL COMMENT '执行者',
    `input`      json     DEFAULT NULL COMMENT '函数入参',
    `output`     json     DEFAULT NULL COMMENT '执行结果',
    `log`        text COMMENT '执行日志',
    `cost_time`  int(8) DEFAULT NULL COMMENT '执行耗时',
    `version`    int(8) NOT NULL COMMENT '版本号',
    `status`     varchar(128) NOT NULL COMMENT '当前状态',
    `created_at` datetime     NOT NULL COMMENT '创建时间',
    `updated_at` datetime     NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
