CREATE TABLE IF NOT EXISTS `finance_exchange_rate` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `from_currency` varchar(10) NOT NULL COMMENT '原币种',
  `to_currency` varchar(10) NOT NULL COMMENT '目标币种',
  `rate` decimal(18,8) NOT NULL COMMENT '汇率（1 原币 = rate 目标币）',
  `source_type` enum('MANUAL') NOT NULL COMMENT '来源类型',
  `source_version` varchar(32) NOT NULL DEFAULT 'v1' COMMENT '来源版本',
  `effective_at` datetime NOT NULL COMMENT '生效时间',
  `status` enum('ACTIVE','INACTIVE') NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
  `remark` text DEFAULT NULL COMMENT '备注',
  `created_by` bigint unsigned NOT NULL COMMENT '创建人',
  `updated_by` bigint unsigned NOT NULL COMMENT '更新人',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `idx_fx_rate_pair_effective` (`from_currency`,`to_currency`,`status`,`effective_at`),
  KEY `idx_fx_rate_effective` (`effective_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='财务汇率表';
