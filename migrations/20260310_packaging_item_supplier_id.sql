ALTER TABLE `packaging_item`
  ADD COLUMN `supplier_id` bigint unsigned DEFAULT NULL COMMENT '供应商ID' AFTER `reorder_quantity`,
  ADD KEY `idx_packaging_item_supplier_id` (`supplier_id`);

UPDATE `packaging_item` pi
INNER JOIN `supplier` s
  ON s.`name` = pi.`supplier_name`
SET pi.`supplier_id` = s.`id`
WHERE pi.`supplier_id` IS NULL
  AND pi.`supplier_name` IS NOT NULL
  AND pi.`supplier_name` <> '';
