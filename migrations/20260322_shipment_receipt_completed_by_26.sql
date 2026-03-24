ALTER TABLE shipment
  ADD COLUMN receipt_completed_by BIGINT UNSIGNED NULL COMMENT '接收完成人' AFTER receipt_completed_at;

