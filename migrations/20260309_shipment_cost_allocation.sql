ALTER TABLE `shipment`
  ADD COLUMN `base_currency` varchar(10) DEFAULT NULL COMMENT '本位币种' AFTER `currency`,
  ADD COLUMN `shipping_cost_base_amount` decimal(18,6) NOT NULL DEFAULT 0 COMMENT '运费本位币金额' AFTER `base_currency`,
  ADD COLUMN `shipping_cost_fx_rate` decimal(18,8) NOT NULL DEFAULT 0 COMMENT '运费汇率' AFTER `shipping_cost_base_amount`,
  ADD COLUMN `shipping_cost_fx_source` varchar(50) DEFAULT NULL COMMENT '运费汇率来源' AFTER `shipping_cost_fx_rate`,
  ADD COLUMN `shipping_cost_fx_version` varchar(32) DEFAULT NULL COMMENT '运费汇率版本' AFTER `shipping_cost_fx_source`,
  ADD COLUMN `shipping_cost_fx_time` datetime DEFAULT NULL COMMENT '运费汇率时间' AFTER `shipping_cost_fx_version`;

ALTER TABLE `finance_cost_event`
  MODIFY COLUMN `event_type` enum('PO_ORDERED','PO_SHIPPED','PO_RECEIVED','PO_ADJUST','SHIPMENT_ALLOCATED') NOT NULL COMMENT '成本事件类型',
  ADD COLUMN `shipment_id` bigint unsigned DEFAULT NULL COMMENT '发货单ID' AFTER `purchase_order_item_id`,
  ADD COLUMN `shipment_item_id` bigint unsigned DEFAULT NULL COMMENT '发货单明细ID' AFTER `shipment_id`,
  ADD KEY `idx_finance_cost_event_shipment` (`shipment_id`),
  ADD KEY `idx_finance_cost_event_shipment_item` (`shipment_item_id`);
