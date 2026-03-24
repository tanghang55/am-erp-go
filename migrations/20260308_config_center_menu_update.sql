UPDATE `menu`
SET
  `title` = '配置中心',
  `title_en` = 'Config Center',
  `path` = '/system/config-center',
  `parent_id` = 1,
  `sort` = 130,
  `is_hidden` = 0,
  `permission_code` = 'system.manage',
  `status` = 'ACTIVE',
  `gmt_modified` = NOW()
WHERE `code` = 'SYSTEM_CONFIG_CENTER';
