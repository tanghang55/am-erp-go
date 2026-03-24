INSERT INTO `finance_exchange_rate`
(`from_currency`,`to_currency`,`rate`,`source_type`,`source_version`,`effective_at`,`status`,`remark`,`created_by`,`updated_by`,`gmt_create`,`gmt_modified`)
SELECT * FROM (
  SELECT 'CNY' AS `from_currency`,'USD' AS `to_currency`,0.14000000 AS `rate`,'MANUAL' AS `source_type`,'v1' AS `source_version`,'2026-01-01 00:00:00' AS `effective_at`,'ACTIVE' AS `status`,'初始化汇率种子' AS `remark`,1 AS `created_by`,1 AS `updated_by`,NOW() AS `gmt_create`,NOW() AS `gmt_modified`
  UNION ALL SELECT 'JPY','USD',0.00670000,'MANUAL','v1','2026-01-01 00:00:00','ACTIVE','初始化汇率种子',1,1,NOW(),NOW()
  UNION ALL SELECT 'AUD','USD',0.66000000,'MANUAL','v1','2026-01-01 00:00:00','ACTIVE','初始化汇率种子',1,1,NOW(),NOW()
  UNION ALL SELECT 'EUR','USD',1.08000000,'MANUAL','v1','2026-01-01 00:00:00','ACTIVE','初始化汇率种子',1,1,NOW(),NOW()
  UNION ALL SELECT 'GBP','USD',1.27000000,'MANUAL','v1','2026-01-01 00:00:00','ACTIVE','初始化汇率种子',1,1,NOW(),NOW()
  UNION ALL SELECT 'CAD','USD',0.74000000,'MANUAL','v1','2026-01-01 00:00:00','ACTIVE','初始化汇率种子',1,1,NOW(),NOW()
) AS seed
WHERE NOT EXISTS (
  SELECT 1
  FROM `finance_exchange_rate` existing
  WHERE existing.`from_currency` = seed.`from_currency`
    AND existing.`to_currency` = seed.`to_currency`
    AND existing.`effective_at` = seed.`effective_at`
);
