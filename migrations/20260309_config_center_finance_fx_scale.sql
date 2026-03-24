INSERT INTO `config_definition`
(`config_key`,`module_code`,`module_name`,`group_code`,`group_name`,`label`,`description`,`value_type`,`scope_type`,`default_value`,`options_json`,`sort`,`is_active`,`gmt_create`,`gmt_modified`)
VALUES
('finance.exchange_rate_scale','finance','财务配置','finance_base','财务基础','汇率小数位','汇率创建与推导结果统一保留的小数位数','INT','GLOBAL','4',NULL,20,1,NOW(),NOW())
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
('finance.exchange_rate_scale','GLOBAL',0,'4',NULL,NOW(),NOW())
ON DUPLICATE KEY UPDATE
`config_value`=VALUES(`config_value`),
`gmt_modified`=VALUES(`gmt_modified`);
