UPDATE `menu`
SET
  `title` = '成本中心',
  `title_en` = 'Cost Center',
  `path` = '/finance/costing',
  `sort` = 620,
  `is_hidden` = 0,
  `status` = 'ACTIVE',
  `gmt_modified` = NOW()
WHERE `code` = 'FINANCE_COSTING';

UPDATE `menu`
SET
  `title` = '财务总览',
  `title_en` = 'Finance Overview',
  `path` = '/finance/profit',
  `sort` = 630,
  `is_hidden` = 0,
  `status` = 'ACTIVE',
  `gmt_modified` = NOW()
WHERE `code` = 'FINANCE_PROFIT';

UPDATE `menu`
SET
  `is_hidden` = 1,
  `gmt_modified` = NOW()
WHERE `code` = 'FINANCE_PRODUCT_COST';
