ALTER TABLE `inventory_balance`
  ADD COLUMN `pending_shipment_reserved` int unsigned NOT NULL DEFAULT 0 COMMENT '待出锁定库存' AFTER `pending_shipment`;

ALTER TABLE `inventory_lot`
  ADD COLUMN `qty_pending_shipment` int unsigned NOT NULL DEFAULT 0 COMMENT '待出库存数量' AFTER `qty_reserved`,
  ADD COLUMN `qty_pending_shipment_reserved` int unsigned NOT NULL DEFAULT 0 COMMENT '待出锁定库存数量' AFTER `qty_pending_shipment`;

ALTER TABLE `inventory_movement`
  MODIFY COLUMN `movement_type` enum(
    'PURCHASE_RECEIPT','SALES_SHIPMENT','SALES_ALLOCATE','SALES_RELEASE','SALES_SHIP',
    'STOCK_TAKE_ADJUSTMENT','MANUAL_ADJUSTMENT','DAMAGE_WRITE_OFF','RETURN_RECEIPT',
    'TRANSFER_OUT','TRANSFER_IN','PURCHASE_SHIP','WAREHOUSE_RECEIVE','INSPECTION_PASS',
    'INSPECTION_FAIL','ASSEMBLY_COMPLETE','SHIPMENT_ALLOCATE','SHIPMENT_RELEASE',
    'LOGISTICS_SHIP','PLATFORM_RECEIVE','RETURN_INSPECT','SHIPMENT_SHIP'
  ) NOT NULL,
  ADD COLUMN `before_pending_shipment` int unsigned NOT NULL DEFAULT 0 COMMENT '变动前待出库存' AFTER `after_sellable_reserved`,
  ADD COLUMN `after_pending_shipment` int unsigned NOT NULL DEFAULT 0 COMMENT '变动后待出库存' AFTER `before_pending_shipment`,
  ADD COLUMN `before_pending_shipment_reserved` int unsigned NOT NULL DEFAULT 0 COMMENT '变动前待出锁定库存' AFTER `after_pending_shipment`,
  ADD COLUMN `after_pending_shipment_reserved` int unsigned NOT NULL DEFAULT 0 COMMENT '变动后待出锁定库存' AFTER `before_pending_shipment_reserved`,
  ADD COLUMN `before_logistics_in_transit` int unsigned NOT NULL DEFAULT 0 COMMENT '变动前物流在途库存' AFTER `after_pending_shipment_reserved`,
  ADD COLUMN `after_logistics_in_transit` int unsigned NOT NULL DEFAULT 0 COMMENT '变动后物流在途库存' AFTER `before_logistics_in_transit`;
