CREATE TABLE `product_category` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `parent_id` bigint unsigned DEFAULT NULL COMMENT '父级品类ID',
  `category_code` varchar(50) NOT NULL COMMENT '品类编码',
  `category_name` varchar(100) NOT NULL COMMENT '品类名称',
  `level` tinyint unsigned NOT NULL COMMENT '层级',
  `status` enum('ACTIVE','INACTIVE') NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_category_code` (`category_code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_level_status` (`level`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='产品品类';
