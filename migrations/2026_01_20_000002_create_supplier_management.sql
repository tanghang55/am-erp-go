-- ============================================
-- Supplier management tables
-- ============================================
CREATE TABLE IF NOT EXISTS supplier (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'Supplier ID',
    supplier_code VARCHAR(50) NOT NULL COMMENT 'Supplier code',
    name VARCHAR(200) NOT NULL COMMENT 'Supplier name',
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' COMMENT 'Status',
    remark TEXT NULL COMMENT 'Remark',
    gmt_create DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
    gmt_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated time',
    UNIQUE KEY uk_supplier_code (supplier_code),
    KEY idx_supplier_name (name),
    KEY idx_supplier_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Supplier table';

CREATE TABLE IF NOT EXISTS supplier_type (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'Primary key',
    supplier_id BIGINT UNSIGNED NOT NULL COMMENT 'Supplier ID',
    type VARCHAR(20) NOT NULL COMMENT 'Type: PRODUCT/PACKAGING/LOGISTICS',
    gmt_create DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
    gmt_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated time',
    UNIQUE KEY uk_supplier_type (supplier_id, type),
    KEY idx_supplier_type_supplier (supplier_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Supplier types';

CREATE TABLE IF NOT EXISTS supplier_contact (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'Primary key',
    supplier_id BIGINT UNSIGNED NOT NULL COMMENT 'Supplier ID',
    name VARCHAR(100) NOT NULL COMMENT 'Contact name',
    phone VARCHAR(50) NULL COMMENT 'Phone',
    email VARCHAR(100) NULL COMMENT 'Email',
    position VARCHAR(100) NULL COMMENT 'Position',
    is_primary TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Is primary contact',
    gmt_create DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
    gmt_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated time',
    KEY idx_supplier_contact_supplier (supplier_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Supplier contacts';

CREATE TABLE IF NOT EXISTS supplier_account (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'Primary key',
    supplier_id BIGINT UNSIGNED NOT NULL COMMENT 'Supplier ID',
    bank_name VARCHAR(100) NOT NULL COMMENT 'Bank name',
    bank_account VARCHAR(100) NOT NULL COMMENT 'Bank account',
    currency VARCHAR(20) NULL COMMENT 'Currency',
    tax_no VARCHAR(100) NULL COMMENT 'Tax number',
    payment_terms VARCHAR(100) NULL COMMENT 'Payment terms',
    gmt_create DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
    gmt_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated time',
    KEY idx_supplier_account_supplier (supplier_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Supplier accounts';

CREATE TABLE IF NOT EXISTS supplier_tag (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'Primary key',
    supplier_id BIGINT UNSIGNED NOT NULL COMMENT 'Supplier ID',
    tag VARCHAR(100) NOT NULL COMMENT 'Tag',
    gmt_create DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
    gmt_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated time',
    KEY idx_supplier_tag_supplier (supplier_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Supplier tags';
