CREATE TABLE IF NOT EXISTS third_party_refund_sync_state (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  provider VARCHAR(32) NOT NULL COMMENT '平台提供方',
  channel VARCHAR(32) NOT NULL COMMENT '同步通道',
  last_posted_after DATETIME NULL COMMENT '最近成功游标',
  last_sync_started_at DATETIME NULL COMMENT '最近同步开始时间',
  last_sync_finished_at DATETIME NULL COMMENT '最近同步结束时间',
  gmt_create DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  gmt_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_provider_channel (provider, channel)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='第三方退款同步状态表';

CREATE TABLE IF NOT EXISTS third_party_refund_sync_run (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  provider VARCHAR(32) NOT NULL COMMENT '平台提供方',
  channel VARCHAR(32) NOT NULL COMMENT '同步通道',
  trigger_type VARCHAR(16) NOT NULL COMMENT '触发类型',
  status VARCHAR(20) NOT NULL COMMENT '任务状态',
  request_posted_after DATETIME NULL COMMENT '请求游标',
  fetched_refunds INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '拉取退款数',
  imported_refunds INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '入库成功数',
  error_refunds INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '失败数',
  message VARCHAR(500) NULL COMMENT '任务消息',
  started_at DATETIME NOT NULL COMMENT '开始时间',
  finished_at DATETIME NULL COMMENT '结束时间',
  gmt_create DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  gmt_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (id),
  KEY idx_provider_channel_started (provider, channel, started_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='第三方退款同步运行记录';

CREATE TABLE IF NOT EXISTS third_party_refund_event (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  provider VARCHAR(32) NOT NULL COMMENT '平台提供方',
  channel VARCHAR(32) NOT NULL COMMENT '同步通道',
  refund_id VARCHAR(128) NOT NULL COMMENT '平台退款唯一标识',
  order_id VARCHAR(64) NOT NULL COMMENT '平台订单号',
  order_item_id VARCHAR(64) NULL COMMENT '平台订单行号',
  seller_sku VARCHAR(100) NOT NULL COMMENT '平台卖家SKU',
  marketplace VARCHAR(10) NOT NULL DEFAULT '' COMMENT '站点编码',
  product_id BIGINT UNSIGNED NULL COMMENT 'ERP产品ID',
  qty_refunded BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '退款数量',
  refund_amount DECIMAL(18,4) NOT NULL DEFAULT 0 COMMENT '退款金额',
  currency CHAR(3) NOT NULL COMMENT '币种',
  posted_at DATETIME NOT NULL COMMENT '平台退款时间',
  status ENUM('MAPPED','UNMAPPED') NOT NULL COMMENT '映射状态',
  error_message VARCHAR(500) NULL COMMENT '错误信息',
  raw_payload TEXT NULL COMMENT '原始快照',
  gmt_create DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  gmt_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_provider_channel_refund_id (provider, channel, refund_id),
  KEY idx_posted_at (posted_at),
  KEY idx_order_id (order_id),
  KEY idx_product_id (product_id),
  KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='第三方退款事件表';
