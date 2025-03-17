CREATE TABLE `app`
(
    `id`         int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `creator`    varchar(255) NOT NULL COMMENT '创建人',
    `name`       varchar(255) NOT NULL COMMENT '应用名',
    `created_at` datetime     NOT NULL COMMENT '创建时间',
    `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime     NOT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_name` (`name`) USING BTREE COMMENT 'name 索引'
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `run_history`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `def_id`     varchar(256) NOT NULL COMMENT '应用名',
    `name`       varchar(256) NOT NULL COMMENT '函数名称',
    `output`     varchar(256) DEFAULT NULL COMMENT '执行结果',
    `run_timer`  varchar(256) NOT NULL COMMENT '运行时间',
    `cost_time`  int(8) DEFAULT NULL COMMENT '执行耗时',
    `status`     varchar(128) NOT NULL COMMENT '当前状态',
    `created_at` datetime     NOT NULL COMMENT '创建时间',
    `updated_at` datetime     NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `deleted_at` datetime     DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE COMMENT '主键索引',
    UNIQUE KEY `idx_def_timer` (`def_id`,`run_timer`) USING BTREE COMMENT '定时器执行时间索引',
    KEY          `idx_created_at` (`created_at`) COMMENT '创建时间索引',
    KEY          `idx_deleted_at` (`deleted_at`) COMMENT '删除时间索引'
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `timer_def`
(
    `id`                int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `app`               varchar(255) NOT NULL COMMENT '应用名',
    `def_id`            varchar(255) NOT NULL COMMENT '定时器唯一ID',
    `name`              varchar(255) NOT NULL COMMENT '定时器name',
    `creator`           varchar(255) NOT NULL COMMENT '创建人',
    `status`            smallint(255) NOT NULL COMMENT '定时器状态 1未激活 2激活',
    `cron`              varchar(255) NOT NULL COMMENT '定时表达式',
    `notify_type`       smallint(255) NOT NULL COMMENT '1 rpc 2 kafka',
    `notify_rpc_param`  json         DEFAULT NULL COMMENT 'rpc参数',
    `trigger_type`      smallint(255) DEFAULT NULL COMMENT '触发类型 1一次 2持续',
    `end_time`          varchar(255) DEFAULT NULL COMMENT '定时器停止时间',
    `notify_http_param` json         DEFAULT NULL COMMENT 'http 参数',
    `timer_type`        smallint(6) NOT NULL COMMENT '定时器类型 1延时 2定时',
    `delay_time`        varchar(255) DEFAULT NULL COMMENT '延时时间',
    `deleted_at`        datetime     DEFAULT NULL COMMENT '删除时间',
    `created_at`        datetime     NOT NULL COMMENT '创建时间',
    `updated_at`        datetime     DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_type`       int(4) DEFAULT '0' COMMENT '删除类型',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_app` (`app`,`name`) USING BTREE COMMENT 'app name 索引',
    KEY                 `idx_def_id` (`def_id`) USING BTREE COMMENT 'def_id 索引\n'
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

ALTER TABLE timer_def
ADD COLUMN execute_time_limit int(4) AFTER delete_type;

