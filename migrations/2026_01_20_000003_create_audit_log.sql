-- ============================================
-- Audit log table
-- ============================================
CREATE TABLE IF NOT EXISTS audit_log (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Log ID',
    trace_id VARCHAR(64) DEFAULT NULL COMMENT 'Trace ID',
    user_id BIGINT UNSIGNED DEFAULT NULL COMMENT 'Operator user ID',
    username VARCHAR(50) DEFAULT NULL COMMENT 'Operator username',
    module VARCHAR(50) NOT NULL COMMENT 'Module',
    action VARCHAR(100) NOT NULL COMMENT 'Action',
    entity_type VARCHAR(50) DEFAULT NULL COMMENT 'Entity type',
    entity_id VARCHAR(100) DEFAULT NULL COMMENT 'Entity ID',
    changes JSON DEFAULT NULL COMMENT 'Changes JSON',
    ip_address VARCHAR(45) DEFAULT NULL COMMENT 'IP address',
    user_agent TEXT COMMENT 'User Agent',
    gmt_create DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
    gmt_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated time',
    PRIMARY KEY (id),
    KEY idx_trace_id (trace_id),
    KEY idx_user_id (user_id),
    KEY idx_module (module),
    KEY idx_action (action),
    KEY idx_entity_type_entity_id (entity_type, entity_id),
    KEY idx_gmt_create (gmt_create)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Audit log table';
