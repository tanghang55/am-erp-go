-- 运费报价表
CREATE TABLE IF NOT EXISTS shipping_rate (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',

    -- 关联信息
    provider_id BIGINT UNSIGNED NOT NULL COMMENT '物流供应商ID',
    origin_warehouse_id BIGINT UNSIGNED NOT NULL COMMENT '起运仓库ID',
    destination_warehouse_id BIGINT UNSIGNED NOT NULL COMMENT '目的仓库ID',

    -- 运输信息
    transport_mode ENUM('EXPRESS', 'AIR', 'SEA', 'RAIL', 'TRUCK') NOT NULL COMMENT '运输方式',
    service_name VARCHAR(100) DEFAULT NULL COMMENT '服务名称(如DHL Express, UPS Standard等)',

    -- 计费方式
    pricing_method ENUM('PER_KG', 'PER_CBM', 'PER_PACKAGE', 'FIXED') NOT NULL COMMENT '计费方式',

    -- 费率
    base_rate DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '基础费率',
    currency VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币',

    -- 最低收费
    min_charge DECIMAL(10,2) DEFAULT NULL COMMENT '最低收费',

    -- 附加费
    fuel_surcharge_rate DECIMAL(5,2) DEFAULT 0.00 COMMENT '燃油附加费率(%)',
    remote_area_surcharge DECIMAL(10,2) DEFAULT NULL COMMENT '偏远地区附加费',

    -- 重量/体积范围（可选，用于阶梯报价）
    min_weight DECIMAL(10,2) DEFAULT NULL COMMENT '最小重量(kg)',
    max_weight DECIMAL(10,2) DEFAULT NULL COMMENT '最大重量(kg)',
    min_volume DECIMAL(10,4) DEFAULT NULL COMMENT '最小体积(m³)',
    max_volume DECIMAL(10,4) DEFAULT NULL COMMENT '最大体积(m³)',

    -- 时效
    transit_days INT DEFAULT NULL COMMENT '运输时效(天)',

    -- 有效期
    effective_date DATE NOT NULL COMMENT '生效日期',
    expiry_date DATE DEFAULT NULL COMMENT '失效日期',

    -- 状态
    status ENUM('ACTIVE', 'INACTIVE', 'EXPIRED') NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',

    -- 备注
    remark TEXT DEFAULT NULL COMMENT '备注',

    -- 审计字段
    created_by BIGINT UNSIGNED DEFAULT NULL COMMENT '创建人',
    updated_by BIGINT UNSIGNED DEFAULT NULL COMMENT '更新人',
    gmt_create DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',

    INDEX idx_provider (provider_id),
    INDEX idx_origin_warehouse (origin_warehouse_id),
    INDEX idx_destination_warehouse (destination_warehouse_id),
    INDEX idx_transport_mode (transport_mode),
    INDEX idx_effective_date (effective_date),
    INDEX idx_expiry_date (expiry_date),
    INDEX idx_status (status),

    FOREIGN KEY (provider_id) REFERENCES logistics_provider(id) ON DELETE RESTRICT,
    FOREIGN KEY (origin_warehouse_id) REFERENCES warehouse(id) ON DELETE RESTRICT,
    FOREIGN KEY (destination_warehouse_id) REFERENCES warehouse(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='运费报价表';
