INSERT INTO `menu` (
  `title`, `title_en`, `code`, `parent_id`, `path`, `icon`, `component`, `sort`, `is_hidden`, `permission_code`, `status`, `gmt_create`, `gmt_modified`
)
SELECT
  '业务配置', 'Business Config', 'SYSTEM_BUSINESS_CONFIG', 1, '/system/business-config', NULL, NULL, 130, 0, 'system.manage', 'ACTIVE', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `menu` WHERE `code` = 'SYSTEM_BUSINESS_CONFIG'
);
