ALTER TABLE `purchase_order`
  ADD COLUMN `closed_at` datetime NULL DEFAULT NULL COMMENT '关闭时间' AFTER `received_at`;

UPDATE `purchase_order`
SET `closed_at` = `gmt_modified`
WHERE `status` = 'CLOSED' AND `closed_at` IS NULL;
