INSERT INTO `config_definition`
(`config_key`,`module_code`,`module_name`,`group_code`,`group_name`,`label`,`description`,`value_type`,`scope_type`,`default_value`,`options_json`,`sort`,`is_active`,`gmt_create`,`gmt_modified`)
VALUES
('procurement.default_currency','procurement','采购配置','procurement_defaults','采购默认规则','采购默认币种','采购单与采购计划转采购单的默认币种','ENUM','GLOBAL','USD','[\"USD\",\"CNY\",\"JPY\",\"AUD\",\"EUR\",\"GBP\",\"CAD\"]',5,1,NOW(),NOW()),
('logistics.default_currency','logistics','物流配置','logistics_defaults','物流默认规则','物流默认币种','物流报价与货件缺省成本币种','ENUM','GLOBAL','CNY','[\"USD\",\"CNY\",\"JPY\",\"AUD\",\"EUR\",\"GBP\",\"CAD\"]',10,1,NOW(),NOW()),
('packaging.default_currency','packaging','包材配置','packaging_defaults','包材默认规则','包材默认币种','包材采购单在物料未维护币种时的默认币种','ENUM','GLOBAL','CNY','[\"USD\",\"CNY\",\"JPY\",\"AUD\",\"EUR\",\"GBP\",\"CAD\"]',10,1,NOW(),NOW())
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

INSERT INTO `config_value`
(`config_key`,`scope_type`,`scope_ref_id`,`config_value`,`updated_by`,`gmt_create`,`gmt_modified`)
VALUES
('procurement.default_currency','GLOBAL',0,'USD',NULL,NOW(),NOW()),
('logistics.default_currency','GLOBAL',0,'CNY',NULL,NOW(),NOW()),
('packaging.default_currency','GLOBAL',0,'CNY',NULL,NOW(),NOW())
ON DUPLICATE KEY UPDATE
`config_value`=VALUES(`config_value`),
`gmt_modified`=VALUES(`gmt_modified`);
