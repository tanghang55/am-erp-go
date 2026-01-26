-- ============================================
-- 添加打包管理菜单到采购管理模块下
-- 2026-01-22
-- ============================================

-- 步骤1: 确保采购管理菜单存在（如果不存在则创建）
INSERT INTO menu (code, title, title_en, path, parent_id, icon, sort, status, gmt_create, gmt_modified)
SELECT 'procurement', '采购管理', 'Procurement', '/procurement', NULL, 'ShoppingCart', 30, 'ACTIVE', NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM menu WHERE code = 'procurement'
);

-- 步骤2: 添加打包管理菜单作为采购管理的子菜单
INSERT INTO menu (code, title, title_en, path, parent_id, icon, sort, status, gmt_create, gmt_modified)
SELECT 'procurement:assembly', '打包管理', 'Assembly Management', '/procurement/assembly',
    (SELECT id FROM menu WHERE code = 'procurement' LIMIT 1),
    'Tools', 20, 'ACTIVE', NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM menu WHERE code = 'procurement:assembly'
);

-- 步骤3: 为管理员角色授权打包管理菜单
-- 注意：假设管理员角色ID为1，请根据实际情况调整    
INSERT INTO role_menu (role_id, menu_id, gmt_create)
SELECT 1, id, NOW()
FROM menu
WHERE code = 'procurement:assembly'
AND NOT EXISTS (
    SELECT 1 FROM role_menu rm
    WHERE rm.role_id = 1 AND rm.menu_id = (SELECT id FROM menu WHERE code = 'procurement:assembly' LIMIT 1)
);

-- 验证SQL（可选，执行后查看结果）
-- SELECT id, code, title, title_en, path, parent_id, sort FROM menu WHERE code IN ('procurement', 'procurement:assembly');
