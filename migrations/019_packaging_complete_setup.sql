-- 包材模块完整设置（包含表创建和类型修复）
-- 如果已执行过 2026_01_14_000002_create_packaging_tables.sql，只执行ALTER部分

-- ============================================================================
-- 1. 创建包材表（如果不存在）
-- ============================================================================
CREATE TABLE IF NOT EXISTS `packaging_items` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `trace_id` VARCHAR(64) NOT NULL COMMENT '追踪ID',
  `item_code` VARCHAR(50) NOT NULL COMMENT '物料编码',
  `item_name` VARCHAR(100) NOT NULL COMMENT '物料名称',
  `category` VARCHAR(50) NOT NULL COMMENT '类别',
  `specification` VARCHAR(200) NULL COMMENT '规格描述',
  `unit_cost` DECIMAL(10,4) NOT NULL DEFAULT 0.0000 COMMENT '单位成本',
  `currency` VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币',
  `unit` VARCHAR(20) NOT NULL DEFAULT 'PCS' COMMENT '单位',
  `quantity_on_hand` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '库存数量',
  `reorder_point` BIGINT UNSIGNED NULL COMMENT '补货点',
  `reorder_quantity` BIGINT UNSIGNED NULL COMMENT '建议补货数量',
  `supplier_name` VARCHAR(100) NULL COMMENT '供应商名称',
  `supplier_contact` VARCHAR(100) NULL COMMENT '供应商联系方式',
  `status` ENUM('ACTIVE', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
  `notes` TEXT NULL COMMENT '备注',
  `created_by` BIGINT UNSIGNED NOT NULL COMMENT '创建人',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_item_code` (`item_code`),
  INDEX `idx_trace_id` (`trace_id`),
  INDEX `idx_category` (`category`),
  INDEX `idx_status` (`status`),
  INDEX `idx_item_name` (`item_name`),
  INDEX `idx_created_by` (`created_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='包装材料主数据表';

CREATE TABLE IF NOT EXISTS `packaging_ledger` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `trace_id` VARCHAR(64) NOT NULL COMMENT '追踪ID',
  `packaging_item_id` BIGINT UNSIGNED NOT NULL COMMENT '包装物料ID',
  `transaction_type` ENUM('IN', 'OUT', 'ADJUSTMENT') NOT NULL COMMENT '流水类型',
  `quantity` BIGINT NOT NULL COMMENT '数量',
  `unit_cost` DECIMAL(10,4) NOT NULL DEFAULT 0.0000 COMMENT '单位成本',
  `total_cost` DECIMAL(15,2) GENERATED ALWAYS AS (ABS(quantity) * unit_cost) STORED COMMENT '总成本',
  `quantity_before` BIGINT UNSIGNED NOT NULL COMMENT '操作前库存',
  `quantity_after` BIGINT UNSIGNED NOT NULL COMMENT '操作后库存',
  `reference_type` VARCHAR(50) NULL COMMENT '关联单据类型',
  `reference_id` BIGINT UNSIGNED NULL COMMENT '关联单据ID',
  `occurred_at` DATETIME NOT NULL COMMENT '发生日期',
  `notes` TEXT NULL COMMENT '备注',
  `created_by` BIGINT UNSIGNED NOT NULL COMMENT '创建人',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  INDEX `idx_trace_id` (`trace_id`),
  INDEX `idx_packaging_item_id` (`packaging_item_id`),
  INDEX `idx_transaction_type` (`transaction_type`),
  INDEX `idx_occurred_at` (`occurred_at`),
  INDEX `idx_reference` (`reference_type`, `reference_id`),
  INDEX `idx_created_by` (`created_by`),
  FOREIGN KEY (`packaging_item_id`) REFERENCES `packaging_items`(`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='包装材料流水表';

-- ============================================================================
-- 2. 如果表已存在但字段类型错误，执行修复
-- ============================================================================
-- 修复packaging_items表的数量字段类型
ALTER TABLE `packaging_items`
  MODIFY COLUMN `quantity_on_hand` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '库存数量',
  MODIFY COLUMN `reorder_point` BIGINT UNSIGNED NULL COMMENT '补货点',
  MODIFY COLUMN `reorder_quantity` BIGINT UNSIGNED NULL COMMENT '建议补货数量';

-- 修复packaging_ledger表的数量字段类型
ALTER TABLE `packaging_ledger`
  MODIFY COLUMN `quantity` BIGINT NOT NULL COMMENT '数量（正数表示入库，负数表示出库）',
  MODIFY COLUMN `quantity_before` BIGINT UNSIGNED NOT NULL COMMENT '操作前库存',
  MODIFY COLUMN `quantity_after` BIGINT UNSIGNED NOT NULL COMMENT '操作后库存';
