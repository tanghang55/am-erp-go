INSERT INTO `config_definition`
(`config_key`,`module_code`,`module_name`,`group_code`,`group_name`,`label`,`description`,`value_type`,`scope_type`,`default_value`,`options_json`,`sort`,`is_active`,`gmt_create`,`gmt_modified`)
VALUES
('finance.default_currency','finance','财务配置','finance_base','财务基础','本位币','系统财务默认本位币','ENUM','GLOBAL','USD','[\"USD\",\"CNY\",\"JPY\",\"AUD\",\"EUR\",\"GBP\",\"CAD\"]',10,1,NOW(),NOW()),
('sales_import.default_channel','sales_import','订单导入配置','sales_import_defaults','导入默认值','默认渠道','CSV列为空时使用','STRING','GLOBAL','MANUAL',NULL,10,1,NOW(),NOW()),
('sales_import.default_marketplace','sales_import','订单导入配置','sales_import_defaults','导入默认值','默认站点','CSV列为空时使用','STRING','GLOBAL','US',NULL,20,1,NOW(),NOW()),
('procurement.demand_window_days','procurement','采购配置','procurement_defaults','采购默认规则','需求窗口天数','采购计划默认统计窗口','INT','GLOBAL','30',NULL,10,1,NOW(),NOW()),
('procurement.default_lead_time_days','procurement','采购配置','procurement_defaults','采购默认规则','默认采购周期','SKU策略未覆盖时使用','INT','GLOBAL','15',NULL,20,1,NOW(),NOW()),
('procurement.default_safety_days','procurement','采购配置','procurement_defaults','采购默认规则','默认安全天数','SKU策略未覆盖时使用','INT','GLOBAL','7',NULL,30,1,NOW(),NOW()),
('procurement.default_moq','procurement','采购配置','procurement_defaults','采购默认规则','默认MOQ','SKU策略未覆盖时使用','INT','GLOBAL','1',NULL,40,1,NOW(),NOW()),
('procurement.default_order_multiple','procurement','采购配置','procurement_defaults','采购默认规则','默认订购倍数','SKU策略未覆盖时使用','INT','GLOBAL','1',NULL,50,1,NOW(),NOW())
ON DUPLICATE KEY UPDATE
`module_code`=VALUES(`module_code`),
`module_name`=VALUES(`module_name`),
`group_code`=VALUES(`group_code`),
`group_name`=VALUES(`group_name`),
`label`=VALUES(`label`),
`description`=VALUES(`description`),
`value_type`=VALUES(`value_type`),
`scope_type`=VALUES(`scope_type`),
`default_value`=VALUES(`default_value`),
`options_json`=VALUES(`options_json`),
`sort`=VALUES(`sort`),
`is_active`=VALUES(`is_active`),
`gmt_modified`=VALUES(`gmt_modified`);
