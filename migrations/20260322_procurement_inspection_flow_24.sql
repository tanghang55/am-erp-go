ALTER TABLE `product`
  ADD COLUMN `is_inspection_required` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '是否需要质检' AFTER `weight_unit_id`;

ALTER TABLE `purchase_order`
  ADD COLUMN `ordered_by` bigint unsigned NULL COMMENT '下单人ID' AFTER `ordered_at`,
  ADD COLUMN `shipped_by` bigint unsigned NULL COMMENT '发货操作人ID' AFTER `shipped_at`,
  ADD COLUMN `received_by` bigint unsigned NULL COMMENT '收货操作人ID' AFTER `received_at`,
  ADD COLUMN `inspected_at` datetime NULL COMMENT '质检时间' AFTER `received_by`,
  ADD COLUMN `inspected_by` bigint unsigned NULL COMMENT '质检操作人ID' AFTER `inspected_at`,
  ADD COLUMN `completed_by` bigint unsigned NULL COMMENT '完成操作人ID' AFTER `closed_at`,
  ADD COLUMN `is_force_completed` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '是否强制完成' AFTER `completed_by`,
  ADD COLUMN `force_completed_at` datetime NULL COMMENT '强制完成时间' AFTER `is_force_completed`,
  ADD COLUMN `force_completed_by` bigint unsigned NULL COMMENT '强制完成操作人ID' AFTER `force_completed_at`,
  ADD COLUMN `force_complete_reason` text NULL COMMENT '强制完成原因' AFTER `force_completed_by`,
  ADD INDEX `idx_ordered_by` (`ordered_by`),
  ADD INDEX `idx_shipped_by` (`shipped_by`),
  ADD INDEX `idx_received_by` (`received_by`),
  ADD INDEX `idx_inspected_by` (`inspected_by`),
  ADD INDEX `idx_completed_by` (`completed_by`),
  ADD INDEX `idx_force_completed_by` (`force_completed_by`);

ALTER TABLE `purchase_order_item`
  ADD COLUMN `qty_inspection_pass` bigint unsigned NOT NULL DEFAULT 0 COMMENT '质检通过数量' AFTER `qty_received`,
  ADD COLUMN `qty_inspection_fail` bigint unsigned NOT NULL DEFAULT 0 COMMENT '质检失败数量' AFTER `qty_inspection_pass`;

ALTER TABLE `inventory_movement`
  MODIFY COLUMN `movement_type` enum(
    'PURCHASE_RECEIPT',
    'SALES_SHIPMENT',
    'SALES_ALLOCATE',
    'SALES_RELEASE',
    'SALES_SHIP',
    'STOCK_TAKE_ADJUSTMENT',
    'MANUAL_ADJUSTMENT',
    'DAMAGE_WRITE_OFF',
    'RETURN_RECEIPT',
    'TRANSFER_OUT',
    'TRANSFER_IN',
    'PURCHASE_SHIP',
    'WAREHOUSE_RECEIVE',
    'INSPECTION_PASS',
    'INSPECTION_FAIL',
    'INSPECTION_LOSS',
    'ASSEMBLY_CONSUME',
    'ASSEMBLY_COMPLETE',
    'SHIPMENT_ALLOCATE',
    'SHIPMENT_RELEASE',
    'LOGISTICS_SHIP',
    'PLATFORM_RECEIVE',
    'RETURN_INSPECT',
    'SHIPMENT_SHIP'
  ) NOT NULL COMMENT '库存流水类型';
