UPDATE `menu`
SET
  `code` = 'INVENTORY_ADJUSTMENTS',
  `gmt_modified` = NOW()
WHERE `code` = 'INVENTORY_MOVEMENTS_CREATE';
