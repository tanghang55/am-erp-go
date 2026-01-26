-- 创建装箱规格表
CREATE TABLE IF NOT EXISTS `package_spec` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(100) NOT NULL COMMENT '名称',
    `length` DECIMAL(10, 2) NOT NULL DEFAULT 0 COMMENT '长(cm)',
    `width` DECIMAL(10, 2) NOT NULL DEFAULT 0 COMMENT '宽(cm)',
    `height` DECIMAL(10, 2) NOT NULL DEFAULT 0 COMMENT '高(cm)',
    `weight` DECIMAL(10, 2) NOT NULL DEFAULT 0 COMMENT '重量(kg)',
    `remark` VARCHAR(500) NULL COMMENT '备注',
    `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
    `created_by` BIGINT UNSIGNED NULL COMMENT '创建人',
    `updated_by` BIGINT UNSIGNED NULL COMMENT '更新人',
    `gmt_create` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `gmt_modified` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_status` (`status`),
    INDEX `idx_gmt_modified` (`gmt_modified`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='装箱规格';

-- 更新发货单明细表，添加装箱信息字段
ALTER TABLE `shipment_item`
ADD COLUMN `package_spec_id` BIGINT UNSIGNED NULL COMMENT '装箱规格ID' AFTER `quantity_shipped`,
ADD COLUMN `box_quantity` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '装箱数量' AFTER `package_spec_id`,
ADD INDEX `idx_package_spec_id` (`package_spec_id`);

-- 插入装箱规格管理菜单
INSERT INTO `menu` (`title`, `title_en`, `code`, `path`, `component`, `icon`, `parent_id`, `sort_order`, `is_visible`, `is_enabled`, `gmt_create`, `gmt_modified`)
SELECT '装箱规格', 'Package Specs', 'shipping.package-specs', '/shipping/package-specs', 'shipping/views/PackageSpecList', 'Box', m.id, 20, 1, 1, NOW(), NOW()
FROM `menu` m WHERE m.code = 'shipping';
