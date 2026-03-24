INSERT INTO `product_config_item` (`config_type`, `item_code`, `item_name`, `status`, `sort`, `remark`)
SELECT seed.config_type, seed.item_code, seed.item_name, seed.status, seed.sort, seed.remark
FROM (
  SELECT 'SALES_STATUS' AS config_type, 'DRAFT' AS item_code, '草稿' AS item_name, 'ACTIVE' AS status, 10 AS sort, '未开始销售' AS remark
  UNION ALL
  SELECT 'SALES_STATUS', 'ON_SALE', '正常销售', 'ACTIVE', 20, '正常销售中的产品'
  UNION ALL
  SELECT 'SALES_STATUS', 'REPLENISHING', '补货中', 'ACTIVE', 30, '需要持续补货的产品'
  UNION ALL
  SELECT 'SALES_STATUS', 'OFF_SHELF', '下架', 'ACTIVE', 40, '停止销售的产品'
) AS seed
WHERE NOT EXISTS (
  SELECT 1
  FROM `product_config_item` existing
  WHERE existing.config_type = seed.config_type
    AND existing.item_code = seed.item_code
);
