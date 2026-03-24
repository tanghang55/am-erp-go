UPDATE `menu`
SET
  `is_hidden` = 0,
  `gmt_modified` = NOW()
WHERE `code` = 'INVENTORY_ADJUSTMENTS';
