UPDATE `menu`
SET
  `title` = '导入记录',
  `title_en` = 'Import History',
  `gmt_modified` = NOW()
WHERE `code` = 'SALES_ORDER_IMPORT';
