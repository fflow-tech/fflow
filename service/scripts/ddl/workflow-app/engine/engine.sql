CREATE TABLE `history_node_inst`
(
    `id`               bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `namespace`        varchar(256) NOT NULL COMMENT '命名空间',
    `def_id`           bigint(20) unsigned NOT NULL COMMENT '主键ID',
    `def_version`      int(8) NOT NULL COMMENT '流程的版本号',
    `inst_id`          bigint(20) unsigned NOT NULL COMMENT '流程实例ID',
    `ref_name`         varchar(256) NOT NULL COMMENT '节点引用名称',
    `context`          json         NOT NULL COMMENT '节点实例上下文',
    `status`           tinyint(4) NOT NULL COMMENT '节点实例状态 1:scheduled,2:waiting,3:paused,4:running,5:completed,6:failed,7:cancelled,8:timeout',
    `scheduled_at`     datetime     NOT NULL COMMENT '节点开始调度时间',
    `wait_at`          datetime DEFAULT NULL COMMENT '节点等待开始时间',
    `execute_at`       datetime DEFAULT NULL COMMENT '节点执行开始时间',
    `asyn_wait_res_at` datetime DEFAULT NULL COMMENT '异步等待结果开始时间',
    `completed_at`     datetime DEFAULT NULL COMMENT '节点执行结束时间',
    `created_at`       datetime     NOT NULL COMMENT '创建时间',
    `updated_at`       datetime     NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `deleted_at`       datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`, `def_id`) USING BTREE COMMENT '主键索引',
    KEY                `idx_def_id_def_version` (`def_id`,`def_version`) USING BTREE COMMENT '流程定义ID索引',
    KEY                `idx_inst_id_ref_name` (`inst_id`,`ref_name`) USING BTREE COMMENT '流程实例ID节点名称索引',
    KEY                `idx_created_at` (`created_at`) COMMENT '创建时间索引',
    KEY                `idx_deleted_at` (`deleted_at`) COMMENT '删除时间索引'
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 shardkey=def_id;

CREATE TABLE `history_workflow_inst`
(
    `id`           bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `namespace`    varchar(256) NOT NULL COMMENT '命名空间',
    `def_id`       bigint(20) unsigned NOT NULL COMMENT '主键ID',
    `def_version`  int(8) NOT NULL COMMENT '流程的版本号',
    `name`         varchar(256) NOT NULL COMMENT '流程实例名称',
    `creator`      varchar(128) NOT NULL COMMENT '创建人',
    `context`      json         NOT NULL COMMENT '流程实例上下文',
    `status`       tinyint(4) NOT NULL COMMENT '流程实例状态 1:running,2:paused,3:completed,4:failed,5:cancelled,6:timeout',
    `start_at`     datetime     NOT NULL COMMENT '流程执行开始时间',
    `completed_at` datetime DEFAULT NULL COMMENT '流程执行结束时间',
    `created_at`   datetime     NOT NULL COMMENT '创建时间',
    `updated_at`   datetime     NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `deleted_at`   datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`, `def_id`) USING BTREE COMMENT '主键索引',
    KEY            `idx_def_id_def_version` (`def_id`,`def_version`) USING BTREE COMMENT '流程定义ID索引',
    KEY            `idx_created_at` (`created_at`) COMMENT '创建时间索引',
    KEY            `idx_deleted_at` (`deleted_at`) COMMENT '删除时间索引'
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 shardkey=def_id;

CREATE TABLE `node_inst`
(
    `id`               bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `namespace`        varchar(256) NOT NULL COMMENT '命名空间',
    `def_id`           bigint(20) unsigned NOT NULL COMMENT '主键ID',
    `def_version`      int(8) NOT NULL COMMENT '流程的版本号',
    `inst_id`          bigint(20) unsigned NOT NULL COMMENT '流程实例ID',
    `ref_name`         varchar(256) NOT NULL COMMENT '节点引用名称',
    `context`          json         NOT NULL COMMENT '节点实例上下文',
    `status`           tinyint(4) NOT NULL COMMENT '节点实例状态 1:scheduled,2:waiting,3:paused,4:running,5:completed,6:failed,7:cancelled,8:timeout',
    `scheduled_at`     datetime     NOT NULL COMMENT '节点开始调度时间',
    `wait_at`          datetime DEFAULT NULL COMMENT '节点等待开始时间',
    `execute_at`       datetime DEFAULT NULL COMMENT '节点执行开始时间',
    `asyn_wait_res_at` datetime DEFAULT NULL COMMENT '异步等待结果开始时间',
    `completed_at`     datetime DEFAULT NULL COMMENT '节点执行结束时间',
    `created_at`       datetime     NOT NULL COMMENT '创建时间',
    `updated_at`       datetime     NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `deleted_at`       datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`, `def_id`) USING BTREE COMMENT '主键索引',
    KEY                `idx_def_id_def_version` (`def_id`,`def_version`) USING BTREE COMMENT '流程定义ID索引',
    KEY                `idx_inst_id_ref_name` (`inst_id`,`ref_name`) USING BTREE COMMENT '流程实例ID节点名称索引',
    KEY                `idx_created_at` (`created_at`) COMMENT '创建时间索引',
    KEY                `idx_deleted_at` (`deleted_at`) COMMENT '删除时间索引'
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 shardkey=def_id;

CREATE TABLE `trigger`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `namespace`   varchar(256) NOT NULL COMMENT '命名空间',
    `type`        varchar(64)  NOT NULL COMMENT '触发器类型',
    `event`       varchar(128) NOT NULL COMMENT '事件名称',
    `expr`        varchar(256) NOT NULL COMMENT '定时触发器的时间表达式',
    `attribute`   json         NOT NULL COMMENT '触发器属性',
    `level`       tinyint(4) NOT NULL COMMENT '触发器类型 1:流程级别触发器, 2:流程实例级别触发器',
    `def_id`      bigint(20) unsigned NOT NULL COMMENT '主键ID',
    `def_version` int(8) NOT NULL COMMENT '流程的版本号',
    `inst_id`     bigint(20) unsigned DEFAULT NULL COMMENT '流程实例ID',
    `status`      tinyint(4) NOT NULL COMMENT '触发器状态，1:未激活, 2:已激活',
    `created_at`  datetime     NOT NULL COMMENT '创建时间',
    `updated_at`  datetime     NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `deleted_at`  datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`, `def_id`) USING BTREE COMMENT '主键索引',
    KEY           `idx_event` (`event`) USING BTREE COMMENT '事件名称索引',
    KEY           `idx_def_id_def_version` (`def_id`,`def_version`) USING BTREE COMMENT '流程定义ID索引',
    KEY           `idx_inst_id` (`inst_id`) USING BTREE COMMENT '流程实例ID索引',
    KEY           `idx_created_at` (`created_at`) COMMENT '创建时间索引',
    KEY           `idx_deleted_at` (`deleted_at`) COMMENT '删除时间索引'
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 shardkey=def_id;

CREATE TABLE `workflow_def`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `namespace`     varchar(256) NOT NULL COMMENT '命名空间',
    `def_id`        bigint(20) unsigned NOT NULL COMMENT '主键ID',
    `parent_def_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '父流程定义ID',
    `attribute`     json         NOT NULL COMMENT '流程定义的其他属性',
    `version`       int(8) NOT NULL COMMENT '流程的版本号',
    `name`          varchar(256) NOT NULL COMMENT '流程定义名称',
    `def_json`      json         NOT NULL COMMENT '流程定义的内容',
    `creator`       varchar(128) NOT NULL COMMENT '创建人',
    `status`        tinyint(4) NOT NULL COMMENT '流程定义状态，1:未激活, 2:已激活',
    `description`   varchar(1024) DEFAULT NULL COMMENT '流程定义描述',
    `created_at`    datetime     NOT NULL COMMENT '创建时间',
    `updated_at`    datetime     NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `deleted_at`    datetime      DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE COMMENT '主键索引',
    UNIQUE KEY `uk_def_id_version` (`def_id`,`version`) USING BTREE COMMENT '流程定义ID索引',
    KEY             `idx_namespace` (`namespace`) COMMENT '命名空间索引',
    KEY             `idx_created_at` (`created_at`) COMMENT '创建时间索引',
    KEY             `idx_deleted_at` (`deleted_at`) COMMENT '删除时间索引'
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `workflow_inst`
(
    `id`           bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `namespace`    varchar(256) NOT NULL COMMENT '命名空间',
    `def_id`       bigint(20) unsigned NOT NULL COMMENT '主键ID',
    `def_version`  int(8) NOT NULL COMMENT '流程的版本号',
    `name`         varchar(256) NOT NULL COMMENT '流程实例名称',
    `creator`      varchar(128) NOT NULL COMMENT '创建人',
    `context`      json         NOT NULL COMMENT '流程实例上下文',
    `status`       tinyint(4) NOT NULL COMMENT '流程实例状态 1:running,2:paused,3:completed,4:failed,5:cancelled,6:timeout',
    `start_at`     datetime     NOT NULL COMMENT '流程执行开始时间',
    `completed_at` datetime DEFAULT NULL COMMENT '流程执行结束时间',
    `created_at`   datetime     NOT NULL COMMENT '创建时间',
    `updated_at`   datetime     NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `deleted_at`   datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`, `def_id`) USING BTREE COMMENT '主键索引',
    KEY            `idx_def_id_def_version` (`def_id`,`def_version`) USING BTREE COMMENT '流程定义ID索引',
    KEY            `idx_created_at` (`created_at`) COMMENT '创建时间索引',
    KEY            `idx_deleted_at` (`deleted_at`) COMMENT '删除时间索引'
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 shardkey=def_id;
