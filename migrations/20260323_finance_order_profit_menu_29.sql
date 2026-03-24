UPDATE `menu`
SET
  `title` = '财务总览',
  `title_en` = 'Finance Overview',
  `path` = '/finance/profit',
  `sort` = 610,
  `is_hidden` = 0,
  `status` = 'ACTIVE',
  `gmt_modified` = NOW()
WHERE `code` = 'FINANCE_PROFIT';

UPDATE `menu`
SET
  `title` = '现金流水',
  `title_en` = 'Cash Ledger',
  `path` = '/finance/cash-ledger',
  `sort` = 620,
  `is_hidden` = 0,
  `status` = 'ACTIVE',
  `gmt_modified` = NOW()
WHERE `code` = 'FINANCE_CASH';

UPDATE `menu`
SET
  `title` = '成本中心',
  `title_en` = 'Cost Center',
  `path` = '/finance/costing',
  `sort` = 630,
  `is_hidden` = 0,
  `status` = 'ACTIVE',
  `gmt_modified` = NOW()
WHERE `code` = 'FINANCE_COSTING';

INSERT INTO `menu` (
  `id`,
  `title`,
  `title_en`,
  `code`,
  `parent_id`,
  `path`,
  `icon`,
  `component`,
  `sort`,
  `is_hidden`,
  `permission_code`,
  `status`,
  `gmt_create`,
  `gmt_modified`
)
SELECT
  57,
  '订单利润',
  'Order Profit',
  'FINANCE_ORDER_PROFIT',
  6,
  '/finance/order-profit',
  NULL,
  NULL,
  640,
  0,
  'finance.manage',
  'ACTIVE',
  NOW(),
  NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `menu` WHERE `code` = 'FINANCE_ORDER_PROFIT'
);

UPDATE `menu`
SET
  `title` = '订单利润',
  `title_en` = 'Order Profit',
  `parent_id` = 6,
  `path` = '/finance/order-profit',
  `sort` = 640,
  `is_hidden` = 0,
  `permission_code` = 'finance.manage',
  `status` = 'ACTIVE',
  `gmt_modified` = NOW()
WHERE `code` = 'FINANCE_ORDER_PROFIT';

UPDATE `menu`
SET
  `title` = '汇率管理',
  `title_en` = 'Exchange Rates',
  `path` = '/finance/exchange-rates',
  `sort` = 650,
  `is_hidden` = 0,
  `status` = 'ACTIVE',
  `gmt_modified` = NOW()
WHERE `code` = 'FINANCE_EXCHANGE_RATES';

