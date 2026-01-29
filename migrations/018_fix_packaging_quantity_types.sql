-- 修复包材表中数量字段的类型
-- 数量字段应该用整数，而不是 DECIMAL

ALTER TABLE `packaging_items`
  MODIFY COLUMN `quantity_on_hand` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '库存数量',
  MODIFY COLUMN `reorder_point` BIGINT UNSIGNED NULL COMMENT '补货点',
  MODIFY COLUMN `reorder_quantity` BIGINT UNSIGNED NULL COMMENT '建议补货数量';

ALTER TABLE `packaging_ledger`
  MODIFY COLUMN `quantity` BIGINT NOT NULL COMMENT '数量（正数表示入库，负数表示出库）',
  MODIFY COLUMN `quantity_before` BIGINT UNSIGNED NOT NULL COMMENT '操作前库存',
  MODIFY COLUMN `quantity_after` BIGINT UNSIGNED NOT NULL COMMENT '操作后库存';
