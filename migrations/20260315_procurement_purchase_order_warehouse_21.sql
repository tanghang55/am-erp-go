ALTER TABLE `purchase_order`
  ADD COLUMN `warehouse_id` bigint unsigned DEFAULT NULL COMMENT '目标仓库ID' AFTER `supplier_id`;
