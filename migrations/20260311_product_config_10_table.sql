CREATE TABLE `product_config_item` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `config_type` varchar(30) NOT NULL COMMENT '配置类型',
  `item_code` varchar(50) NOT NULL COMMENT '配置编码',
  `item_name` varchar(100) NOT NULL COMMENT '配置名称',
  `status` enum('ACTIVE','INACTIVE') NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_product_config_type_code` (`config_type`,`item_code`),
  UNIQUE KEY `uk_product_config_type_name` (`config_type`,`item_name`),
  KEY `idx_product_config_type_status` (`config_type`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='产品配置项';
