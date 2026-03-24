CREATE TABLE IF NOT EXISTS `config_value` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `config_key` varchar(120) NOT NULL COMMENT '配置键',
  `scope_type` varchar(20) NOT NULL DEFAULT 'GLOBAL' COMMENT '作用域类型',
  `scope_ref_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '作用域引用ID，GLOBAL固定0',
  `config_value` varchar(255) NOT NULL COMMENT '配置值',
  `updated_by` bigint unsigned DEFAULT NULL COMMENT '更新人',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_value_scope` (`config_key`,`scope_type`,`scope_ref_id`),
  KEY `idx_config_value_scope` (`scope_type`,`scope_ref_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='配置中心值表';
