INSERT INTO `config_value`
(`config_key`,`scope_type`,`scope_ref_id`,`config_value`,`updated_by`,`gmt_create`,`gmt_modified`)
SELECT
  seed.`config_key`,
  'GLOBAL',
  0,
  seed.`config_value`,
  seed.`updated_by`,
  NOW(),
  NOW()
FROM (
  SELECT
    CONVERT('finance.default_currency' USING utf8mb4) COLLATE utf8mb4_unicode_ci AS `config_key`,
    CONVERT(UPPER(TRIM(COALESCE((SELECT `config_value` FROM `business_config` WHERE `config_key` = 'finance.default_currency' ORDER BY `id` DESC LIMIT 1), 'USD'))) USING utf8mb4) COLLATE utf8mb4_unicode_ci AS `config_value`,
    (SELECT `updated_by` FROM `business_config` WHERE `config_key` = 'finance.default_currency' ORDER BY `id` DESC LIMIT 1) AS `updated_by`
  UNION ALL
  SELECT
    CONVERT('sales_import.default_channel' USING utf8mb4) COLLATE utf8mb4_unicode_ci AS `config_key`,
    CONVERT(TRIM(COALESCE((SELECT `config_value` FROM `business_config` WHERE `config_key` = 'sales_import.default_channel' ORDER BY `id` DESC LIMIT 1), 'MANUAL')) USING utf8mb4) COLLATE utf8mb4_unicode_ci AS `config_value`,
    (SELECT `updated_by` FROM `business_config` WHERE `config_key` = 'sales_import.default_channel' ORDER BY `id` DESC LIMIT 1) AS `updated_by`
  UNION ALL
  SELECT
    CONVERT('sales_import.default_marketplace' USING utf8mb4) COLLATE utf8mb4_unicode_ci AS `config_key`,
    CONVERT(UPPER(TRIM(COALESCE((SELECT `config_value` FROM `business_config` WHERE `config_key` = 'sales_import.default_marketplace' ORDER BY `id` DESC LIMIT 1), 'US'))) USING utf8mb4) COLLATE utf8mb4_unicode_ci AS `config_value`,
    (SELECT `updated_by` FROM `business_config` WHERE `config_key` = 'sales_import.default_marketplace' ORDER BY `id` DESC LIMIT 1) AS `updated_by`
  UNION ALL
  SELECT
    CONVERT('procurement.demand_window_days' USING utf8mb4) COLLATE utf8mb4_unicode_ci AS `config_key`,
    CONVERT(CAST(COALESCE((SELECT `demand_window_days` FROM `procurement_replenishment_config` ORDER BY `id` ASC LIMIT 1), 30) AS CHAR) USING utf8mb4) COLLATE utf8mb4_unicode_ci AS `config_value`,
    NULL AS `updated_by`
  UNION ALL
  SELECT
    CONVERT('procurement.default_lead_time_days' USING utf8mb4) COLLATE utf8mb4_unicode_ci AS `config_key`,
    CONVERT(CAST(COALESCE((SELECT `default_lead_time_days` FROM `procurement_replenishment_config` ORDER BY `id` ASC LIMIT 1), 15) AS CHAR) USING utf8mb4) COLLATE utf8mb4_unicode_ci AS `config_value`,
    NULL AS `updated_by`
  UNION ALL
  SELECT
    CONVERT('procurement.default_safety_days' USING utf8mb4) COLLATE utf8mb4_unicode_ci AS `config_key`,
    CONVERT(CAST(COALESCE((SELECT `default_safety_days` FROM `procurement_replenishment_config` ORDER BY `id` ASC LIMIT 1), 7) AS CHAR) USING utf8mb4) COLLATE utf8mb4_unicode_ci AS `config_value`,
    NULL AS `updated_by`
  UNION ALL
  SELECT
    CONVERT('procurement.default_moq' USING utf8mb4) COLLATE utf8mb4_unicode_ci AS `config_key`,
    CONVERT(CAST(COALESCE((SELECT `default_moq` FROM `procurement_replenishment_config` ORDER BY `id` ASC LIMIT 1), 1) AS CHAR) USING utf8mb4) COLLATE utf8mb4_unicode_ci AS `config_value`,
    NULL AS `updated_by`
  UNION ALL
  SELECT
    CONVERT('procurement.default_order_multiple' USING utf8mb4) COLLATE utf8mb4_unicode_ci AS `config_key`,
    CONVERT(CAST(COALESCE((SELECT `default_order_multiple` FROM `procurement_replenishment_config` ORDER BY `id` ASC LIMIT 1), 1) AS CHAR) USING utf8mb4) COLLATE utf8mb4_unicode_ci AS `config_value`,
    NULL AS `updated_by`
) AS seed
ON DUPLICATE KEY UPDATE
`config_value`=VALUES(`config_value`),
`updated_by`=VALUES(`updated_by`),
`gmt_modified`=VALUES(`gmt_modified`);
