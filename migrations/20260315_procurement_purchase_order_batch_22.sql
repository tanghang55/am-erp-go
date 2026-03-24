ALTER TABLE `purchase_order`
  ADD COLUMN `batch_no` varchar(50) NOT NULL DEFAULT '' COMMENT '采购批次号' AFTER `po_number`;

UPDATE `purchase_order`
SET `batch_no` = `po_number`
WHERE (`batch_no` IS NULL OR `batch_no` = '')
  AND `po_number` IS NOT NULL
  AND `po_number` <> '';

ALTER TABLE `purchase_order`
  ADD KEY `idx_batch_no` (`batch_no`);
