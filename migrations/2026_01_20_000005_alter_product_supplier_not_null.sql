-- ============================================
-- 产品默认供应商必填
-- ============================================
ALTER TABLE product MODIFY supplier_id BIGINT UNSIGNED NOT NULL COMMENT '默认供应商ID';
