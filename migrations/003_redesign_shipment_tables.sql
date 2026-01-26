-- 删除旧的发货单表
DROP TABLE IF EXISTS `shipment_item`;
DROP TABLE IF EXISTS `shipment`;

-- 重新创建发货单表 - Amazon业务场景
CREATE TABLE IF NOT EXISTS `shipment` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `shipment_number` VARCHAR(50) NOT NULL UNIQUE COMMENT '内部发货单号',

    -- 订单信息
    `order_number` VARCHAR(100) NULL COMMENT '订单号',
    `sales_channel` VARCHAR(50) NULL COMMENT '销售渠道: Amazon/eBay/Shopify等',
    `marketplace` VARCHAR(10) NULL COMMENT '站点: US/UK/DE/JP等',

    -- Amazon FBA相关字段
    `shipment_type` ENUM('FBA', 'FBM', 'WHOLESALE', 'INTERNAL') NOT NULL DEFAULT 'FBA' COMMENT '发货类型',
    `shipment_plan_id` VARCHAR(50) NULL COMMENT 'Amazon Shipment Plan ID',
    `amazon_shipment_id` VARCHAR(50) NULL COMMENT 'Amazon分配的Shipment ID',
    `amazon_reference_id` VARCHAR(50) NULL COMMENT 'Amazon Reference ID',
    `destination_fc_id` VARCHAR(20) NULL COMMENT '目标FBA仓库代码 (PHX3/ONT8等)',
    `fc_address` TEXT NULL COMMENT 'FBA仓库地址',
    `label_prep_type` VARCHAR(50) NULL COMMENT '贴标方式: SELLER_LABEL/AMAZON_LABEL',

    -- 仓库和物流信息
    `warehouse_id` BIGINT UNSIGNED NOT NULL COMMENT '发货仓库ID',
    `carrier` VARCHAR(50) NULL COMMENT '承运商: FedEx/UPS/DHL等',
    `shipping_method` VARCHAR(50) NULL COMMENT '运输方式: 海运/空运/快递/陆运',
    `tracking_number` VARCHAR(100) NULL COMMENT '物流追踪号',
    `carrier_tracking_url` VARCHAR(500) NULL COMMENT '承运商追踪链接',

    -- 包装信息
    `box_count` INT UNSIGNED DEFAULT 0 COMMENT '箱数',
    `total_weight` DECIMAL(10,2) DEFAULT 0 COMMENT '总重量(kg)',
    `total_volume` DECIMAL(10,3) DEFAULT 0 COMMENT '总体积(m³)',
    `pallet_count` INT UNSIGNED DEFAULT 0 COMMENT '托盘数',

    -- 费用信息
    `shipping_cost` DECIMAL(12,4) NOT NULL DEFAULT 0 COMMENT '运费',
    `prep_cost` DECIMAL(12,4) NOT NULL DEFAULT 0 COMMENT '预处理费用',
    `other_cost` DECIMAL(12,4) NOT NULL DEFAULT 0 COMMENT '其他费用',
    `currency` VARCHAR(10) NOT NULL DEFAULT 'USD' COMMENT '币种',

    -- 时间节点
    `expected_ship_date` DATE NULL COMMENT '预计发货日期',
    `expected_arrival_date` DATE NULL COMMENT '预计到达日期',
    `actual_ship_date` DATETIME NULL COMMENT '实际发货时间',
    `actual_arrival_date` DATETIME NULL COMMENT '实际到达时间',
    `checked_in_at` DATETIME NULL COMMENT 'Amazon签收时间',

    -- 状态管理
    `status` ENUM(
        'DRAFT',           -- 草稿
        'WORKING',         -- 工作中(拣货/打包)
        'READY_TO_SHIP',   -- 待发货
        'SHIPPED',         -- 已发货
        'IN_TRANSIT',      -- 运输中
        'RECEIVING',       -- 收货中(Amazon正在收货)
        'RECEIVED',        -- 已收货(Amazon收货完成)
        'CHECKED_IN',      -- 已签收(全部上架)
        'CLOSED',          -- 已关闭
        'CANCELLED'        -- 已取消
    ) NOT NULL DEFAULT 'DRAFT' COMMENT '发货状态',

    `inventory_locked` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '库存是否已锁定',
    `inventory_shipped` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '库存是否已出库',

    -- 备注和附件
    `remark` TEXT NULL COMMENT '备注',
    `prep_instructions` TEXT NULL COMMENT '预处理说明',
    `attachment_urls` TEXT NULL COMMENT '附件URL (JSON数组)',

    -- 审计字段
    `created_by` BIGINT UNSIGNED NULL COMMENT '创建人',
    `updated_by` BIGINT UNSIGNED NULL COMMENT '更新人',
    `gmt_create` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `gmt_modified` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX `idx_order_number` (`order_number`),
    INDEX `idx_shipment_plan_id` (`shipment_plan_id`),
    INDEX `idx_amazon_shipment_id` (`amazon_shipment_id`),
    INDEX `idx_warehouse_id` (`warehouse_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_shipment_type` (`shipment_type`),
    INDEX `idx_tracking_number` (`tracking_number`),
    INDEX `idx_actual_ship_date` (`actual_ship_date`),
    INDEX `idx_gmt_create` (`gmt_create`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='发货单';

-- 发货单明细表
CREATE TABLE IF NOT EXISTS `shipment_item` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `shipment_id` BIGINT UNSIGNED NOT NULL COMMENT '发货单ID',
    `sku_id` BIGINT UNSIGNED NOT NULL COMMENT 'SKU ID',

    -- Amazon相关字段
    `asin` VARCHAR(20) NULL COMMENT 'ASIN',
    `fnsku` VARCHAR(20) NULL COMMENT 'FNSKU (Amazon仓库SKU)',
    `msku` VARCHAR(100) NULL COMMENT 'MSKU (Merchant SKU = seller_sku)',

    -- 数量相关
    `quantity_planned` INT UNSIGNED NOT NULL COMMENT '计划发货数量',
    `quantity_shipped` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '实际发货数量',
    `quantity_received` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Amazon收货数量',
    `quantity_in_case` INT UNSIGNED NULL COMMENT '每箱数量',

    -- 预处理信息
    `prep_type` VARCHAR(50) NULL COMMENT '预处理类型: Polybagging/Bubble/Taping/Labeling等',
    `prep_owner` VARCHAR(50) NULL COMMENT '预处理负责方: AMAZON/SELLER',

    -- 成本信息
    `unit_cost` DECIMAL(12,4) NOT NULL DEFAULT 0 COMMENT '单位成本',
    `currency` VARCHAR(10) NOT NULL DEFAULT 'USD' COMMENT '币种',

    -- 库存流转标记
    `inventory_locked` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '库存已锁定',
    `inventory_deducted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '库存已扣减',

    `gmt_create` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `gmt_modified` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX `idx_shipment_id` (`shipment_id`),
    INDEX `idx_sku_id` (`sku_id`),
    INDEX `idx_asin` (`asin`),
    INDEX `idx_fnsku` (`fnsku`),
    FOREIGN KEY (`shipment_id`) REFERENCES `shipment`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='发货单明细';

-- 发货单箱唛信息表
CREATE TABLE IF NOT EXISTS `shipment_box` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `shipment_id` BIGINT UNSIGNED NOT NULL COMMENT '发货单ID',
    `box_number` VARCHAR(50) NOT NULL COMMENT '箱号',
    `box_label` VARCHAR(100) NULL COMMENT '箱唛/Box Label',
    `tracking_id` VARCHAR(100) NULL COMMENT '箱子的追踪号',
    `weight` DECIMAL(10,2) DEFAULT 0 COMMENT '重量(kg)',
    `length` DECIMAL(10,2) DEFAULT 0 COMMENT '长(cm)',
    `width` DECIMAL(10,2) DEFAULT 0 COMMENT '宽(cm)',
    `height` DECIMAL(10,2) DEFAULT 0 COMMENT '高(cm)',
    `volume` DECIMAL(10,3) DEFAULT 0 COMMENT '体积(m³)',
    `gmt_create` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `gmt_modified` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY `uk_shipment_box` (`shipment_id`, `box_number`),
    FOREIGN KEY (`shipment_id`) REFERENCES `shipment`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='发货单箱唛信息';

-- 发货单箱内商品明细表
CREATE TABLE IF NOT EXISTS `shipment_box_item` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `box_id` BIGINT UNSIGNED NOT NULL COMMENT '箱子ID',
    `shipment_item_id` BIGINT UNSIGNED NOT NULL COMMENT '发货单明细ID',
    `quantity` INT UNSIGNED NOT NULL COMMENT '该箱内的数量',
    `gmt_create` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX `idx_box_id` (`box_id`),
    INDEX `idx_shipment_item_id` (`shipment_item_id`),
    FOREIGN KEY (`box_id`) REFERENCES `shipment_box`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`shipment_item_id`) REFERENCES `shipment_item`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='发货单箱内商品明细';

-- 发货单状态变更日志表
CREATE TABLE IF NOT EXISTS `shipment_status_log` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `shipment_id` BIGINT UNSIGNED NOT NULL COMMENT '发货单ID',
    `from_status` VARCHAR(20) NULL COMMENT '原状态',
    `to_status` VARCHAR(20) NOT NULL COMMENT '新状态',
    `remark` TEXT NULL COMMENT '备注',
    `operator_id` BIGINT UNSIGNED NULL COMMENT '操作人ID',
    `gmt_create` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX `idx_shipment_id` (`shipment_id`),
    INDEX `idx_gmt_create` (`gmt_create`),
    FOREIGN KEY (`shipment_id`) REFERENCES `shipment`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='发货单状态变更日志';
