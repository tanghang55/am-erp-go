UPDATE `menu`
SET
  `path` = '/inventory/adjustments',
  `gmt_modified` = NOW()
WHERE `code` = 'INVENTORY_MOVEMENTS_CREATE';
