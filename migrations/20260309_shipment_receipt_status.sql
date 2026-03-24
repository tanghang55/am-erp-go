ALTER TABLE `shipment`
  ADD COLUMN `receipt_status` enum('PENDING','PARTIAL','COMPLETED') NOT NULL DEFAULT 'PENDING' COMMENT '接收状态' AFTER `status`,
  ADD COLUMN `receipt_completed_at` datetime DEFAULT NULL COMMENT '接收完成时间' AFTER `delivered_at`;
