ALTER TABLE `inventory_lot`
  ADD COLUMN `qty_purchasing_in_transit` int unsigned NOT NULL DEFAULT 0 COMMENT '采购在途数量' AFTER `qty_in`,
  ADD COLUMN `qty_pending_inspection` int unsigned NOT NULL DEFAULT 0 COMMENT '待检数量' AFTER `qty_purchasing_in_transit`,
  ADD COLUMN `qty_raw_material` int unsigned NOT NULL DEFAULT 0 COMMENT '原料数量' AFTER `qty_pending_inspection`;

ALTER TABLE `inventory_movement`
  ADD COLUMN `before_purchasing_in_transit` int unsigned NOT NULL DEFAULT 0 COMMENT '变动前采购在途' AFTER `after_reserved`,
  ADD COLUMN `after_purchasing_in_transit` int unsigned NOT NULL DEFAULT 0 COMMENT '变动后采购在途' AFTER `before_purchasing_in_transit`,
  ADD COLUMN `before_pending_inspection` int unsigned NOT NULL DEFAULT 0 COMMENT '变动前待检' AFTER `after_purchasing_in_transit`,
  ADD COLUMN `after_pending_inspection` int unsigned NOT NULL DEFAULT 0 COMMENT '变动后待检' AFTER `before_pending_inspection`,
  ADD COLUMN `before_raw_material` int unsigned NOT NULL DEFAULT 0 COMMENT '变动前原料' AFTER `after_pending_inspection`,
  ADD COLUMN `after_raw_material` int unsigned NOT NULL DEFAULT 0 COMMENT '变动后原料' AFTER `before_raw_material`;
