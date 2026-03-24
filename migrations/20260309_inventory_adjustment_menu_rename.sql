UPDATE `menu`
SET
  `title` = '调整库存',
  `title_en` = 'Adjust Inventory',
  `gmt_modified` = NOW()
WHERE `code` = 'INVENTORY_MOVEMENTS_CREATE';
