ALTER TABLE `finance_cost_event`
  MODIFY COLUMN `event_type` enum('PO_ORDERED','PO_SHIPPED','PO_RECEIVED','PO_ADJUST','PACKING_MATERIAL','SHIPMENT_ALLOCATED') NOT NULL COMMENT '成本事件类型',
  ADD COLUMN `inventory_movement_id` bigint unsigned DEFAULT NULL COMMENT '库存流水ID' AFTER `shipment_item_id`,
  ADD KEY `idx_finance_cost_event_inventory_movement` (`inventory_movement_id`);
