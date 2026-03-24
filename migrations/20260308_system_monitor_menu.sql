INSERT INTO menu (
  title,
  title_en,
  code,
  parent_id,
  path,
  component,
  icon,
  sort,
  is_hidden,
  permission_code,
  status
)
SELECT
  '系统监控',
  'System Monitor',
  'SYSTEM_MONITOR',
  NULL,
  '/system/monitor',
  NULL,
  NULL,
  101,
  0,
  'system.manage',
  'ACTIVE'
FROM dual
WHERE NOT EXISTS (
  SELECT 1 FROM menu WHERE code = 'SYSTEM_MONITOR'
);
