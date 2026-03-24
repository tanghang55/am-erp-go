CREATE TABLE IF NOT EXISTS `procurement_replenishment_plan_purchase_order` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `plan_id` bigint unsigned NOT NULL COMMENT '补货计划ID',
  `purchase_order_id` bigint unsigned NOT NULL COMMENT '采购单ID',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_plan_purchase_order` (`plan_id`, `purchase_order_id`),
  KEY `idx_purchase_order_id` (`purchase_order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='补货计划与采购单关联表';

INSERT IGNORE INTO `procurement_replenishment_plan_purchase_order` (`plan_id`, `purchase_order_id`)
SELECT `id`, `purchase_order_id`
FROM `procurement_replenishment_plan`
WHERE `purchase_order_id` IS NOT NULL;
