-- ============================================
-- 产品供应商报价表
-- ============================================
CREATE TABLE IF NOT EXISTS product_supplier_quote (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '报价ID',
    product_id BIGINT UNSIGNED NOT NULL COMMENT '产品ID',
    supplier_id BIGINT UNSIGNED NOT NULL COMMENT '供应商ID',
    price DECIMAL(15,4) NOT NULL COMMENT '报价',
    currency CHAR(3) NOT NULL COMMENT '币种',
    qty_moq INT UNSIGNED NOT NULL DEFAULT 1 COMMENT '起订量',
    lead_time_days INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '交期(天)',
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
    remark TEXT NULL COMMENT '备注',
    gmt_create DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_product_supplier (product_id, supplier_id),
    KEY idx_product_id (product_id),
    KEY idx_supplier_id (supplier_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='产品供应商报价表';
