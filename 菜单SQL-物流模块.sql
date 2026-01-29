-- ============================================
-- 物流管理模块菜单 SQL
-- ============================================
-- 说明：请根据实际情况选择以下方案之一执行

-- ============================================
-- 方案一：使用自动ID（推荐）
-- ============================================
-- 此方案不需要手动指定ID，由数据库自动分配

-- Step 1: 添加物流管理一级菜单
INSERT INTO menus (parent_id, name, path, component, icon, sort, hidden, status) VALUES
(0, '物流管理', '/logistics', '', 'Truck', 60, 0, 1);

-- Step 2: 记录刚插入的父菜单ID（在MySQL客户端中执行）
-- 执行后记下返回的ID，假设为 51，则在后续语句中替换 {PARENT_ID}

-- Step 3: 添加二级菜单（将 {PARENT_ID} 替换为实际的父菜单ID）
INSERT INTO menus (parent_id, name, path, component, icon, sort, hidden, status) VALUES
({PARENT_ID}, '物流供应商', '/logistics/providers', 'logistics/providers/index', 'Shop', 1, 0, 1),
({PARENT_ID}, '运费报价', '/logistics/shipping-rates', 'logistics/shipping-rates/index', 'Money', 2, 0, 1);

-- Step 4: 为管理员角色分配菜单权限（将 {MENU_ID} 替换为实际的菜单ID）
INSERT INTO role_menus (role_id, menu_id) VALUES
(1, {PARENT_ID}),      -- 物流管理
(1, {MENU_ID_1}),      -- 物流供应商
(1, {MENU_ID_2});      -- 运费报价


-- ============================================
-- 方案二：一次性执行（推荐用于MySQL 8.0+）
-- ============================================
-- 使用事务和变量，一次性完成所有操作

START TRANSACTION;

-- 插入一级菜单
INSERT INTO menus (parent_id, name, path, component, icon, sort, hidden, status) VALUES
(0, '物流管理', '/logistics', '', 'Truck', 60, 0, 1);

-- 获取刚插入的ID
SET @parent_menu_id = LAST_INSERT_ID();

-- 插入二级菜单
INSERT INTO menus (parent_id, name, path, component, icon, sort, hidden, status) VALUES
(@parent_menu_id, '物流供应商', '/logistics/providers', 'logistics/providers/index', 'Shop', 1, 0, 1),
(@parent_menu_id, '运费报价', '/logistics/shipping-rates', 'logistics/shipping-rates/index', 'Money', 2, 0, 1);

-- 获取二级菜单ID
SET @provider_menu_id = LAST_INSERT_ID() - 1;
SET @rate_menu_id = LAST_INSERT_ID();

-- 分配给管理员角色（角色ID为1）
INSERT INTO role_menus (role_id, menu_id) VALUES
(1, @parent_menu_id),
(1, @provider_menu_id),
(1, @rate_menu_id);

COMMIT;


-- ============================================
-- 方案三：手动指定ID（需要先查询当前最大ID）
-- ============================================
-- 执行前请先查询：SELECT MAX(id) FROM menus;
-- 假设最大ID为50，则从51开始

INSERT INTO menus (id, parent_id, name, path, component, icon, sort, hidden, status) VALUES
(51, 0, '物流管理', '/logistics', '', 'Truck', 60, 0, 1),
(52, 51, '物流供应商', '/logistics/providers', 'logistics/providers/index', 'Shop', 1, 0, 1),
(53, 51, '运费报价', '/logistics/shipping-rates', 'logistics/shipping-rates/index', 'Money', 2, 0, 1);

-- 分配给管理员角色
INSERT INTO role_menus (role_id, menu_id) VALUES
(1, 51),
(1, 52),
(1, 53);


-- ============================================
-- 验证查询
-- ============================================
-- 执行后用以下SQL验证是否插入成功

-- 查看物流管理菜单
SELECT * FROM menus WHERE name LIKE '%物流%' OR path LIKE '%logistics%';

-- 查看管理员角色的物流菜单权限
SELECT m.* FROM menus m
JOIN role_menus rm ON m.id = rm.menu_id
WHERE rm.role_id = 1 AND (m.name LIKE '%物流%' OR m.path LIKE '%logistics%');


-- ============================================
-- 菜单字段说明
-- ============================================
-- parent_id: 父菜单ID，0表示一级菜单
-- name: 菜单显示名称
-- path: 前端路由路径
-- component: 前端组件路径（一级菜单为空）
-- icon: Element Plus 图标名称
--   - Truck: 卡车图标（物流）
--   - Shop: 商店图标（供应商）
--   - Money: 金钱图标（报价）
-- sort: 排序，数字越小越靠前
-- hidden: 是否隐藏，0=显示，1=隐藏
-- status: 状态，1=启用，0=禁用


-- ============================================
-- 删除菜单（如需重新执行）
-- ============================================
-- 警告：删除菜单会级联删除角色菜单关联，谨慎执行！

-- DELETE FROM menus WHERE name = '物流管理';
-- DELETE FROM menus WHERE parent_id IN (SELECT id FROM menus WHERE name = '物流管理');
