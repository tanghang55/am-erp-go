INSERT INTO `menu`
(`title`,`title_en`,`code`,`parent_id`,`path`,`icon`,`component`,`sort`,`is_hidden`,`permission_code`,`status`,`gmt_create`,`gmt_modified`)
SELECT
  '配置中心', 'Config Center', 'SYSTEM_CONFIG_CENTER', 1, '/system/config-center', NULL, NULL, 130, 0, 'system.manage', 'ACTIVE', NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1 FROM `menu` WHERE `code` = 'SYSTEM_CONFIG_CENTER'
);
