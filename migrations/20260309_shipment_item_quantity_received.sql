ALTER TABLE `shipment_item`
  ADD COLUMN `quantity_received` int unsigned NOT NULL DEFAULT 0 COMMENT '实际接收数量' AFTER `quantity_shipped`;
