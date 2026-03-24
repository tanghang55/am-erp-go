ALTER TABLE `sales_order`
  ADD COLUMN `stock_pool` enum('AVAILABLE','SELLABLE') NOT NULL DEFAULT 'AVAILABLE' COMMENT '库存池' AFTER `marketplace`;

ALTER TABLE `inventory_balance`
  ADD COLUMN `sellable_reserved` int unsigned NOT NULL DEFAULT 0 COMMENT '可售锁定库存' AFTER `reserved_quantity`;

ALTER TABLE `inventory_movement`
  ADD COLUMN `stock_pool` enum('AVAILABLE','SELLABLE') DEFAULT NULL COMMENT '库存池' AFTER `reference_number`,
  ADD COLUMN `before_sellable` int unsigned NOT NULL DEFAULT 0 COMMENT '变动前可售库存' AFTER `after_reserved`,
  ADD COLUMN `after_sellable` int unsigned NOT NULL DEFAULT 0 COMMENT '变动后可售库存' AFTER `before_sellable`,
  ADD COLUMN `before_sellable_reserved` int unsigned NOT NULL DEFAULT 0 COMMENT '变动前可售锁定库存' AFTER `after_sellable`,
  ADD COLUMN `after_sellable_reserved` int unsigned NOT NULL DEFAULT 0 COMMENT '变动后可售锁定库存' AFTER `before_sellable_reserved`,
  ADD KEY `idx_inventory_movement_stock_pool` (`stock_pool`);

ALTER TABLE `inventory_lot`
  ADD COLUMN `qty_sellable` int unsigned NOT NULL DEFAULT 0 COMMENT '可售数量' AFTER `qty_reserved`,
  ADD COLUMN `qty_sellable_reserved` int unsigned NOT NULL DEFAULT 0 COMMENT '可售锁定数量' AFTER `qty_sellable`;
