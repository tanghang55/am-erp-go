-- Users table
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    real_name VARCHAR(100) DEFAULT '' COMMENT '真实姓名',
    email VARCHAR(100) DEFAULT '' COMMENT '邮箱',
    phone VARCHAR(20) DEFAULT '' COMMENT '手机号',
    avatar VARCHAR(255) DEFAULT '' COMMENT '头像URL',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 1=启用, 0=禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- Roles table
CREATE TABLE IF NOT EXISTS roles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE COMMENT '角色名',
    code VARCHAR(50) NOT NULL UNIQUE COMMENT '角色代码',
    description VARCHAR(255) DEFAULT '' COMMENT '描述',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 1=启用, 0=禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_code (code),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- User-Roles junction table
CREATE TABLE IF NOT EXISTS user_roles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    role_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_role (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- Permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '权限名',
    code VARCHAR(100) NOT NULL UNIQUE COMMENT '权限代码',
    description VARCHAR(255) DEFAULT '' COMMENT '描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- Role-Permissions junction table
CREATE TABLE IF NOT EXISTS role_permissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    role_id BIGINT UNSIGNED NOT NULL,
    permission_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_role_permission (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- Menus table
CREATE TABLE IF NOT EXISTS menus (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    parent_id BIGINT UNSIGNED DEFAULT 0 COMMENT '父菜单ID',
    name VARCHAR(50) NOT NULL COMMENT '菜单名',
    path VARCHAR(255) DEFAULT '' COMMENT '路由路径',
    component VARCHAR(255) DEFAULT '' COMMENT '组件路径',
    icon VARCHAR(50) DEFAULT '' COMMENT '图标',
    sort INT DEFAULT 0 COMMENT '排序',
    hidden TINYINT NOT NULL DEFAULT 0 COMMENT '是否隐藏: 1=隐藏, 0=显示',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 1=启用, 0=禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_parent_id (parent_id),
    INDEX idx_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';

-- Role-Menus junction table
CREATE TABLE IF NOT EXISTS role_menus (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    role_id BIGINT UNSIGNED NOT NULL,
    menu_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_role_menu (role_id, menu_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色菜单关联表';

-- ============================================
-- Initial Data
-- ============================================

-- Insert admin role
INSERT INTO roles (name, code, description, status) VALUES
('超级管理员', 'admin', '系统超级管理员，拥有所有权限', 1);

-- Insert admin user (password: admin123)
-- bcrypt hash for 'admin123': $2a$10$N9qo8uLOickgx2ZMRZoMy.MQDsMCYjI8oq8/dJX0Zn5VPm5k5h2VO
INSERT INTO users (username, password_hash, real_name, email, status) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MQDsMCYjI8oq8/dJX0Zn5VPm5k5h2VO', '系统管理员', 'admin@example.com', 1);

-- Link admin user to admin role
INSERT INTO user_roles (user_id, role_id) VALUES (1, 1);

-- Insert default permissions
INSERT INTO permissions (name, code, description) VALUES
('系统管理', 'system:manage', '系统管理权限'),
('用户管理', 'user:manage', '用户管理权限'),
('角色管理', 'role:manage', '角色管理权限'),
('菜单管理', 'menu:manage', '菜单管理权限');

-- Link admin role to all permissions
INSERT INTO role_permissions (role_id, permission_id) VALUES
(1, 1), (1, 2), (1, 3), (1, 4);

-- Insert default menus
INSERT INTO menus (parent_id, name, path, component, icon, sort, hidden, status) VALUES
(0, '首页', '/dashboard', 'dashboard/index', 'HomeFilled', 1, 0, 1),
(0, '系统管理', '/system', '', 'Setting', 100, 0, 1),
(2, '用户管理', '/system/users', 'system/users/index', 'User', 1, 0, 1),
(2, '角色管理', '/system/roles', 'system/roles/index', 'UserFilled', 2, 0, 1),
(2, '菜单管理', '/system/menus', 'system/menus/index', 'Menu', 3, 0, 1);

-- Link admin role to all menus
INSERT INTO role_menus (role_id, menu_id) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5);
