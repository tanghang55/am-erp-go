SET @sql = IF(
  EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'integration_authorization'
      AND column_name = 'refresh_fail_count'
  ),
  'SELECT 1',
  'ALTER TABLE integration_authorization ADD COLUMN refresh_fail_count TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT ''连续刷新失败次数'' AFTER last_refresh_at'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql = IF(
  EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'integration_authorization'
      AND column_name = 'last_refresh_attempt_at'
  ),
  'SELECT 1',
  'ALTER TABLE integration_authorization ADD COLUMN last_refresh_attempt_at DATETIME NULL COMMENT ''最近一次刷新尝试时间'' AFTER refresh_fail_count'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql = IF(
  EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'integration_authorization'
      AND column_name = 'last_refresh_failed_at'
  ),
  'SELECT 1',
  'ALTER TABLE integration_authorization ADD COLUMN last_refresh_failed_at DATETIME NULL COMMENT ''最近一次刷新失败时间'' AFTER last_refresh_attempt_at'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS integration_sku_mapping (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  provider_code VARCHAR(64) NOT NULL COMMENT '平台提供方编码',
  marketplace CHAR(10) NOT NULL COMMENT '站点编码',
  seller_sku VARCHAR(100) NOT NULL COMMENT '平台卖家SKU',
  product_id BIGINT UNSIGNED NOT NULL COMMENT 'ERP产品ID',
  status ENUM('ACTIVE','DISABLED') NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
  remark VARCHAR(255) NULL COMMENT '备注',
  created_by BIGINT UNSIGNED NULL COMMENT '创建人',
  updated_by BIGINT UNSIGNED NULL COMMENT '更新人',
  gmt_create DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  gmt_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_provider_marketplace_seller_sku (provider_code, marketplace, seller_sku),
  KEY idx_product_id (product_id),
  KEY idx_status (status),
  KEY idx_provider_marketplace_status (provider_code, marketplace, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='平台SKU与ERP产品映射表';
