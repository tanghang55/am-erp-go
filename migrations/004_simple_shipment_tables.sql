-- 删除之前的复杂设计
DROP TABLE IF EXISTS `shipment_status_log`;
DROP TABLE IF EXISTS `shipment_box_item`;
DROP TABLE IF EXISTS `shipment_box`;
DROP TABLE IF EXISTS `shipment_item`;
DROP TABLE IF EXISTS `shipment`;

-- 发货单主表 - 通用设计，淡化平台
CREATE TABLE IF NOT EXISTS `shipment` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `shipment_number` VARCHAR(50) NOT NULL UNIQUE COMMENT '发货单号',

    -- 订单信息（通用）
    `order_number` VARCHAR(100) NULL COMMENT '订单号/Reference',
    `sales_channel` VARCHAR(50) NULL COMMENT '销售渠道: Amazon/eBay/Shopify/Offline等',

    -- 发货仓库
    `warehouse_id` BIGINT UNSIGNED NOT NULL COMMENT '发货仓库ID',

    -- 收货方信息（通用）
    `destination_type` ENUM('PLATFORM_WAREHOUSE', 'CUSTOMER', 'OWN_WAREHOUSE', 'SUPPLIER', 'OTHER')
        NOT NULL DEFAULT 'PLATFORM_WAREHOUSE' COMMENT '收货方类型',
    `destination_name` VARCHAR(200) NULL COMMENT '收货方名称',
    `destination_contact` VARCHAR(100) NULL COMMENT '收货联系人',
    `destination_phone` VARCHAR(50) NULL COMMENT '收货电话',
    `destination_address` TEXT NULL COMMENT '收货地址',
    `destination_code` VARCHAR(50) NULL COMMENT '收货方代码(如FBA仓库代码)',

    -- 物流信息（通用）
    `carrier` VARCHAR(50) NULL COMMENT '承运商',
    `shipping_method` VARCHAR(50) NULL COMMENT '运输方式: 快递/空运/海运/陆运',
    `tracking_number` VARCHAR(200) NULL COMMENT '物流追踪号(可能多个,逗号分隔)',

    -- 包装信息
    `box_count` INT UNSIGNED DEFAULT 0 COMMENT '箱数',
    `total_weight` DECIMAL(10,2) DEFAULT 0 COMMENT '总重量(kg)',
    `total_volume` DECIMAL(10,3) DEFAULT 0 COMMENT '总体积(m³)',

    -- 费用
    `shipping_cost` DECIMAL(12,4) NOT NULL DEFAULT 0 COMMENT '运费',
    `currency` VARCHAR(10) NOT NULL DEFAULT 'USD' COMMENT '币种',

    -- 时间节点
    `ship_date` DATE NULL COMMENT '发货日期',
    `expected_delivery_date` DATE NULL COMMENT '预计到达日期',
    `actual_delivery_date` DATE NULL COMMENT '实际到达日期',

    -- 状态（简化）
    `status` ENUM(
        'DRAFT',        -- 草稿
        'CONFIRMED',    -- 已确认(库存已锁定)
        'PACKED',       -- 已打包(库存已扣减到待出)
        'SHIPPED',      -- 已发货(库存已扣减到在途)
        'DELIVERED',    -- 已送达
        'CANCELLED'     -- 已取消
    ) NOT NULL DEFAULT 'DRAFT' COMMENT '状态',

    -- 库存标记
    `inventory_locked` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '库存是否已锁定',
    `inventory_deducted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '库存是否已扣减',

    -- 备注
    `remark` TEXT NULL COMMENT '备注',
    `internal_notes` TEXT NULL COMMENT '内部备注',

    -- 审计
    `created_by` BIGINT UNSIGNED NULL COMMENT '创建人',
    `updated_by` BIGINT UNSIGNED NULL COMMENT '更新人',
    `gmt_create` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `gmt_modified` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX `idx_order_number` (`order_number`),
    INDEX `idx_warehouse_id` (`warehouse_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_destination_type` (`destination_type`),
    INDEX `idx_tracking_number` (`tracking_number`),
    INDEX `idx_ship_date` (`ship_date`),
    INDEX `idx_gmt_create` (`gmt_create`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='发货单';

-- 发货单明细表
CREATE TABLE IF NOT EXISTS `shipment_item` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `shipment_id` BIGINT UNSIGNED NOT NULL COMMENT '发货单ID',
    `sku_id` BIGINT UNSIGNED NOT NULL COMMENT 'SKU ID',

    -- 数量
    `quantity_planned` INT UNSIGNED NOT NULL COMMENT '计划发货数量',
    `quantity_shipped` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '实际发货数量',

    -- 成本
    `unit_cost` DECIMAL(12,4) NOT NULL DEFAULT 0 COMMENT '单位成本',
    `currency` VARCHAR(10) NOT NULL DEFAULT 'USD' COMMENT '币种',

    -- 备注（如果某个SKU有特殊说明）
    `remark` VARCHAR(500) NULL COMMENT '备注',

    `gmt_create` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `gmt_modified` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX `idx_shipment_id` (`shipment_id`),
    INDEX `idx_sku_id` (`sku_id`),
    FOREIGN KEY (`shipment_id`) REFERENCES `shipment`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='发货单明细';

-- 库存流水增加新的流转类型
ALTER TABLE `inventory_movement`
MODIFY COLUMN `movement_type` ENUM(
    'PURCHASE_RECEIPT',
    'SALES_SHIPMENT',
    'STOCK_TAKE_ADJUSTMENT',
    'MANUAL_ADJUSTMENT',
    'DAMAGE_WRITE_OFF',
    'RETURN_RECEIPT',
    'TRANSFER_OUT',
    'TRANSFER_IN',
    'PURCHASE_SHIP',
    'WAREHOUSE_RECEIVE',
    'INSPECTION_PASS',
    'INSPECTION_FAIL',
    'ASSEMBLY_COMPLETE',
    'LOGISTICS_SHIP',
    'PLATFORM_RECEIVE',
    'RETURN_INSPECT',
    'SHIPMENT_LOCK',        -- 发货单锁定库存
    'SHIPMENT_PACK',        -- 发货单打包(原料→待出)
    'SHIPMENT_SHIP',        -- 发货单发货(待出→在途)
    'SHIPMENT_ROLLBACK'     -- 发货单回滚
) NOT NULL COMMENT '流水类型';
