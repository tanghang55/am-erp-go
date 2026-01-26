-- Create shipment table
CREATE TABLE IF NOT EXISTS `shipment` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `shipment_number` VARCHAR(50) NOT NULL UNIQUE COMMENT '发货单号',
    `order_number` VARCHAR(100) NULL COMMENT '订单号',
    `warehouse_id` BIGINT UNSIGNED NOT NULL COMMENT '仓库ID',
    `status` ENUM('PENDING', 'PROCESSING', 'SHIPPED', 'DELIVERED', 'CANCELLED') NOT NULL DEFAULT 'PENDING' COMMENT '状态',
    `carrier` VARCHAR(50) NULL COMMENT '承运商',
    `tracking_number` VARCHAR(100) NULL COMMENT '物流单号',
    `shipping_cost` DECIMAL(12,4) NOT NULL DEFAULT 0 COMMENT '运费',
    `currency` VARCHAR(10) NOT NULL DEFAULT 'USD' COMMENT '币种',
    `shipped_at` DATETIME NULL COMMENT '发货时间',
    `delivered_at` DATETIME NULL COMMENT '送达时间',
    `remark` TEXT NULL COMMENT '备注',
    `created_by` BIGINT UNSIGNED NULL COMMENT '创建人',
    `updated_by` BIGINT UNSIGNED NULL COMMENT '更新人',
    `gmt_create` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `gmt_modified` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_warehouse_id` (`warehouse_id`),
    INDEX `idx_order_number` (`order_number`),
    INDEX `idx_status` (`status`),
    INDEX `idx_tracking_number` (`tracking_number`),
    INDEX `idx_shipped_at` (`shipped_at`),
    INDEX `idx_gmt_create` (`gmt_create`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='发货单';

-- Create shipment_item table
CREATE TABLE IF NOT EXISTS `shipment_item` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `shipment_id` BIGINT UNSIGNED NOT NULL COMMENT '发货单ID',
    `sku_id` BIGINT UNSIGNED NOT NULL COMMENT 'SKU ID',
    `quantity` INT UNSIGNED NOT NULL COMMENT '发货数量',
    `unit_cost` DECIMAL(12,4) NOT NULL DEFAULT 0 COMMENT '单位成本',
    `currency` VARCHAR(10) NOT NULL DEFAULT 'USD' COMMENT '币种',
    `gmt_create` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `gmt_modified` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_shipment_id` (`shipment_id`),
    INDEX `idx_sku_id` (`sku_id`),
    FOREIGN KEY (`shipment_id`) REFERENCES `shipment`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='发货单明细';

-- Update inventory_movement table enum to include new movement types
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
    'RETURN_INSPECT'
) NOT NULL COMMENT '流水类型';
