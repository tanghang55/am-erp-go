-- 产品包材关联表
-- 用于维护每个产品需要消耗的包材及数量

CREATE TABLE IF NOT EXISTS product_packaging_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',

    -- 关联
    product_id BIGINT UNSIGNED NOT NULL COMMENT '产品ID (关联 sku 表)',
    packaging_item_id BIGINT UNSIGNED NOT NULL COMMENT '包材ID (关联 packaging_item 表)',

    -- 消耗数量
    quantity_per_unit DECIMAL(10,3) NOT NULL COMMENT '每个产品需要的包材数量',

    -- 备注
    notes VARCHAR(500) DEFAULT NULL COMMENT '备注说明',

    -- 审计字段
    created_by BIGINT UNSIGNED DEFAULT NULL COMMENT '创建人',
    updated_by BIGINT UNSIGNED DEFAULT NULL COMMENT '更新人',
    gmt_create DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',

    -- 索引
    INDEX idx_product_id (product_id),
    INDEX idx_packaging_item_id (packaging_item_id),
    UNIQUE KEY uk_product_packaging (product_id, packaging_item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品包材关联表';
