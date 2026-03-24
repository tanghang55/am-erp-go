ALTER TABLE `product`
  ADD COLUMN `brand_id` bigint unsigned DEFAULT NULL COMMENT '品牌ID' AFTER `supplier_id`,
  ADD COLUMN `category_id` bigint unsigned DEFAULT NULL COMMENT '品类ID' AFTER `brand_id`,
  ADD COLUMN `sales_unit_id` bigint unsigned DEFAULT NULL COMMENT '销售单位ID' AFTER `category_id`,
  ADD COLUMN `dimension_unit_id` bigint unsigned DEFAULT NULL COMMENT '尺寸单位ID' AFTER `sales_unit_id`,
  ADD COLUMN `weight_unit_id` bigint unsigned DEFAULT NULL COMMENT '重量单位ID' AFTER `dimension_unit_id`,
  ADD COLUMN `length` decimal(10,2) DEFAULT NULL COMMENT '长度' AFTER `weight`,
  ADD COLUMN `width` decimal(10,2) DEFAULT NULL COMMENT '宽度' AFTER `length`,
  ADD COLUMN `height` decimal(10,2) DEFAULT NULL COMMENT '高度' AFTER `width`,
  ADD KEY `idx_brand_id` (`brand_id`),
  ADD KEY `idx_category_id` (`category_id`),
  ADD KEY `idx_sales_unit_id` (`sales_unit_id`),
  ADD KEY `idx_dimension_unit_id` (`dimension_unit_id`),
  ADD KEY `idx_weight_unit_id` (`weight_unit_id`);
