INSERT INTO `product_config_item` (`config_type`, `item_code`, `item_name`, `status`, `sort`, `remark`)
VALUES
  ('SALES_UNIT', 'PCS', '件', 'ACTIVE', 10, '默认销售单位'),
  ('SALES_UNIT', 'SET', '套', 'ACTIVE', 20, '组合或套装单位'),
  ('SALES_UNIT', 'BOX', '箱', 'ACTIVE', 30, '整箱销售单位'),
  ('DIMENSION_UNIT', 'CM', '厘米', 'ACTIVE', 10, '默认尺寸单位'),
  ('DIMENSION_UNIT', 'M', '米', 'ACTIVE', 20, '大件尺寸单位'),
  ('DIMENSION_UNIT', 'IN', '英寸', 'ACTIVE', 30, '英制尺寸单位'),
  ('WEIGHT_UNIT', 'KG', '千克', 'ACTIVE', 10, '默认重量单位'),
  ('WEIGHT_UNIT', 'G', '克', 'ACTIVE', 20, '轻小件重量单位'),
  ('WEIGHT_UNIT', 'LB', '磅', 'ACTIVE', 30, '英制重量单位');
