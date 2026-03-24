INSERT INTO `menu`
(`title`,`title_en`,`code`,`parent_id`,`path`,`icon`,`component`,`sort`,`is_hidden`,`permission_code`,`status`,`gmt_create`,`gmt_modified`)
SELECT
  '汇率管理',
  'Exchange Rates',
  'FINANCE_EXCHANGE_RATES',
  parent.`id`,
  '/finance/exchange-rates',
  'Coin',
  NULL,
  650,
  0,
  'finance.manage',
  'ACTIVE',
  NOW(),
  NOW()
FROM `menu` parent
WHERE parent.`code` = 'FINANCE'
  AND NOT EXISTS (
    SELECT 1 FROM `menu` m WHERE m.`code` = 'FINANCE_EXCHANGE_RATES'
  );
