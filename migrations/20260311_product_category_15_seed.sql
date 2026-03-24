INSERT INTO `product_category` (`parent_id`, `category_code`, `category_name`, `level`, `status`, `sort`, `remark`)
VALUES
  (NULL, 'FINISHED', '成品', 1, 'ACTIVE', 10, '固定一级品类'),
  (NULL, 'SEMI', '半成品', 1, 'ACTIVE', 20, '固定一级品类'),
  (NULL, 'RAW', '原材料', 1, 'ACTIVE', 30, '固定一级品类'),
  (NULL, 'ACCESSORY', '配件', 1, 'ACTIVE', 40, '固定一级品类'),
  (NULL, 'OTHER', '其他', 1, 'ACTIVE', 50, '固定一级品类');
