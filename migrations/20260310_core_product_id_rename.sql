ALTER TABLE `inventory_balance`
  RENAME COLUMN `sku_id` TO `product_id`;

ALTER TABLE `inventory_lot`
  RENAME COLUMN `sku_id` TO `product_id`;

ALTER TABLE `inventory_movement`
  RENAME COLUMN `sku_id` TO `product_id`;

ALTER TABLE `purchase_order_item`
  RENAME COLUMN `sku_id` TO `product_id`;

ALTER TABLE `procurement_replenishment_policy`
  RENAME COLUMN `sku_id` TO `product_id`;

ALTER TABLE `procurement_replenishment_item`
  RENAME COLUMN `sku_id` TO `product_id`;

ALTER TABLE `procurement_replenishment_strategy`
  RENAME COLUMN `sku_id` TO `product_id`;

ALTER TABLE `procurement_replenishment_plan`
  RENAME COLUMN `sku_id` TO `product_id`;

ALTER TABLE `sales_order_item`
  RENAME COLUMN `sku_id` TO `product_id`;

ALTER TABLE `costing_snapshot`
  RENAME COLUMN `sku_id` TO `product_id`;

ALTER TABLE `finance_cost_event`
  RENAME COLUMN `sku_id` TO `product_id`;

ALTER TABLE `finance_order_cost_detail`
  RENAME COLUMN `sku_id` TO `product_id`;

ALTER TABLE `shipment_item`
  RENAME COLUMN `sku_id` TO `product_id`;
