CREATE TABLE `product_combo` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '组合记录ID',
  `combo_id` bigint unsigned NOT NULL COMMENT '组合ID',
  `main_product_id` bigint unsigned NOT NULL COMMENT '主产品ID',
  `product_id` bigint unsigned NOT NULL COMMENT '组件产品ID',
  `qty_ratio` int unsigned NOT NULL DEFAULT '1' COMMENT '数量比例',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_combo_product` (`combo_id`,`product_id`),
  KEY `idx_combo_id` (`combo_id`),
  KEY `idx_main_product_id` (`main_product_id`),
  KEY `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='产品组合表';
