-- ============================================
-- 产品图片表（多图 + 排序）
-- 遵循 erp-schema 规范
-- ============================================
CREATE TABLE IF NOT EXISTS product_image (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    product_id BIGINT UNSIGNED NOT NULL COMMENT '产品ID',
    image_url VARCHAR(500) NOT NULL COMMENT '图片URL',
    sort_order INT UNSIGNED NOT NULL DEFAULT 1 COMMENT '排序序号',
    is_primary TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否主图(1是/0否)',
    status VARCHAR(20) NULL COMMENT '状态',
    remark TEXT NULL COMMENT '备注',
    gmt_create DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_product_image (product_id, image_url),
    KEY idx_product_id_sort (product_id, sort_order),
    KEY idx_product_id_primary (product_id, is_primary)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='产品图片表';
