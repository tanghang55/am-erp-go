DELETE FROM `config_value`
WHERE `config_key` IN (
  'procurement.default_currency',
  'logistics.default_currency',
  'packaging.default_currency'
);

DELETE FROM `config_definition`
WHERE `config_key` IN (
  'procurement.default_currency',
  'logistics.default_currency',
  'packaging.default_currency'
);
