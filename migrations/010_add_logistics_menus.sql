-- ============================================
-- 物流管理菜单
-- ============================================

-- 添加物流管理一级菜单
INSERT INTO menus (parent_id, name, path, component, icon, sort, hidden, status) VALUES
(0, '物流管理', '/logistics', '', 'Truck', 60, 0, 1);

-- 获取刚插入的物流管理菜单ID（假设为 @logistics_parent_id）
SET @logistics_parent_id = LAST_INSERT_ID();

-- 添加物流管理二级菜单
INSERT INTO menus (parent_id, name, path, component, icon, sort, hidden, status) VALUES
(@logistics_parent_id, '物流供应商', '/logistics/providers', 'logistics/providers/index', 'Shop', 1, 0, 1),
(@logistics_parent_id, '运费报价', '/logistics/shipping-rates', 'logistics/shipping-rates/index', 'Money', 2, 0, 1);

-- 为管理员角色分配物流菜单权限（假设管理员角色ID为1）
-- 获取新增菜单的ID
SET @provider_menu_id = LAST_INSERT_ID() - 1;
SET @rate_menu_id = LAST_INSERT_ID();

-- 分配给管理员角色
INSERT INTO role_menus (role_id, menu_id) VALUES
(1, @logistics_parent_id),
(1, @provider_menu_id),
(1, @rate_menu_id);

-- 说明：
-- icon 使用 Element Plus 的图标名称
-- Truck - 物流管理主菜单
-- Shop - 物流供应商
-- Money - 运费报价
-- sort 排序设置为 60，通常在产品、供应商、采购、库存、发货模块之后
