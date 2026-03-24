CREATE TABLE IF NOT EXISTS `business_config` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `config_key` varchar(100) NOT NULL COMMENT '配置键',
  `config_group` varchar(50) NOT NULL COMMENT '配置分组',
  `config_value` varchar(255) NOT NULL COMMENT '配置值',
  `value_type` varchar(20) NOT NULL DEFAULT 'STRING' COMMENT '值类型',
  `description` varchar(255) DEFAULT NULL COMMENT '说明',
  `updated_by` bigint unsigned DEFAULT NULL COMMENT '最后更新人',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_business_config_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='业务配置';
