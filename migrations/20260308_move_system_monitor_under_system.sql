UPDATE menu
SET
  parent_id = 1,
  sort = 150,
  permission_code = 'system.manage'
WHERE code = 'SYSTEM_MONITOR';
