-- ============================================
-- 快速添加打包管理菜单（简化版）
-- 直接在MySQL中执行此SQL
-- ============================================

USE am_erp;

-- 方案1: 如果采购管理菜单已存在，直接添加打包管理子菜单
-- 先查询采购管理的ID（假设code='procurement'或类似）
-- 替换下面的 parent_id 值为实际的采购管理菜单ID

INSERT INTO menu (code, title, title_en, path, parent_id, icon, sort, status, gmt_create, gmt_modified)
VALUES (
    'procurement:assembly',
    '打包管理',
    'Assembly Management',
    '/procurement/assembly',
    (SELECT id FROM (SELECT id FROM menu WHERE code = 'procurement' LIMIT 1) AS tmp),  -- 采购管理菜单ID
    'Tools',
    20,
    'ACTIVE',
    NOW(),
    NOW()
);

-- 为管理员角色授权（假设role_id=1）
INSERT INTO role_menu (role_id, menu_id, gmt_create)
VALUES (
    1,
    (SELECT id FROM menu WHERE code = 'procurement:assembly' LIMIT 1),
    NOW()
);

-- 查询验证
SELECT m.id, m.code, m.title, m.title_en, m.path, m.parent_id, m.sort,
       p.code as parent_code, p.title as parent_title
FROM menu m
LEFT JOIN menu p ON m.parent_id = p.id
WHERE m.code IN ('procurement', 'procurement:assembly')
ORDER BY m.parent_id, m.sort;
