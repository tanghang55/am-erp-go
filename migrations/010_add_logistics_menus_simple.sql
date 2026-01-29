-- ============================================
-- 物流管理菜单（简化版本 - 直接指定ID）
-- ============================================

-- 注意：执行前请先查询当前最大菜单ID，然后相应调整下面的ID值
-- SELECT MAX(id) FROM menus;

-- 假设当前菜单最大ID为50，则从51开始

-- 添加物流管理一级菜单 (ID: 51)
INSERT INTO menus (id, parent_id, name, path, component, icon, sort, hidden, status) VALUES
(51, 0, '物流管理', '/logistics', '', 'Truck', 60, 0, 1);

-- 添加物流管理二级菜单
INSERT INTO menus (id, parent_id, name, path, component, icon, sort, hidden, status) VALUES
(52, 51, '物流供应商', '/logistics/providers', 'logistics/providers/index', 'Shop', 1, 0, 1),
(53, 51, '运费报价', '/logistics/shipping-rates', 'logistics/shipping-rates/index', 'Money', 2, 0, 1);

-- 为管理员角色分配物流菜单权限（假设管理员角色ID为1）
INSERT INTO role_menus (role_id, menu_id) VALUES
(1, 51),
(1, 52),
(1, 53);
