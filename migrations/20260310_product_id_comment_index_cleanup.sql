ALTER TABLE `product`
  MODIFY COLUMN `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '产品ID';

ALTER TABLE `inventory_balance`
  MODIFY COLUMN `product_id` bigint unsigned NOT NULL COMMENT '产品ID',
  RENAME INDEX `uk_sku_warehouse` TO `uk_product_warehouse`,
  RENAME INDEX `idx_sku_id` TO `idx_product_id`;

ALTER TABLE `inventory_lot`
  MODIFY COLUMN `product_id` bigint unsigned NOT NULL COMMENT '产品ID',
  RENAME INDEX `uk_sku_warehouse_lot` TO `uk_product_warehouse_lot`,
  RENAME INDEX `idx_sku_status` TO `idx_product_status`,
  RENAME INDEX `idx_warehouse_sku_received` TO `idx_warehouse_product_received`;

ALTER TABLE `inventory_movement`
  MODIFY COLUMN `product_id` bigint unsigned NOT NULL COMMENT '产品ID',
  RENAME INDEX `idx_sku_id` TO `idx_product_id`,
  RENAME INDEX `idx_sku_warehouse_operated` TO `idx_product_warehouse_operated`;

ALTER TABLE `purchase_order_item`
  MODIFY COLUMN `product_id` bigint unsigned NOT NULL COMMENT '产品ID',
  RENAME INDEX `idx_sku_id` TO `idx_product_id`;

ALTER TABLE `sales_order_item`
  MODIFY COLUMN `product_id` bigint unsigned NOT NULL COMMENT '产品ID',
  RENAME INDEX `idx_sku_id` TO `idx_product_id`;

ALTER TABLE `shipment_item`
  MODIFY COLUMN `product_id` bigint unsigned NOT NULL COMMENT '产品ID',
  RENAME INDEX `idx_sku_id` TO `idx_product_id`;

ALTER TABLE `costing_snapshot`
  MODIFY COLUMN `product_id` bigint unsigned NOT NULL COMMENT '产品ID',
  RENAME INDEX `idx_sku_id` TO `idx_product_id`,
  RENAME INDEX `idx_sku_cost_from` TO `idx_product_cost_from`,
  RENAME INDEX `idx_sku_cost_to` TO `idx_product_cost_to`;

ALTER TABLE `finance_cost_event`
  MODIFY COLUMN `product_id` bigint unsigned NOT NULL COMMENT '产品ID',
  RENAME INDEX `idx_sku_warehouse` TO `idx_product_warehouse`;

ALTER TABLE `finance_order_cost_detail`
  MODIFY COLUMN `product_id` bigint unsigned NOT NULL COMMENT '产品ID',
  RENAME INDEX `idx_sku_warehouse` TO `idx_product_warehouse`;

ALTER TABLE `procurement_replenishment_policy`
  MODIFY COLUMN `product_id` bigint unsigned DEFAULT NULL COMMENT '产品ID',
  RENAME INDEX `uk_sku_id` TO `uk_product_id`;

ALTER TABLE `procurement_replenishment_item`
  MODIFY COLUMN `product_id` bigint unsigned NOT NULL COMMENT '产品ID',
  RENAME INDEX `idx_sku_warehouse` TO `idx_product_warehouse`;

ALTER TABLE `procurement_replenishment_strategy`
  MODIFY COLUMN `product_id` bigint unsigned DEFAULT NULL COMMENT '产品ID',
  RENAME INDEX `idx_sku_id` TO `idx_product_id`;

ALTER TABLE `procurement_replenishment_plan`
  RENAME INDEX `uk_plan_date_sku_wh` TO `uk_plan_date_product_wh`;
