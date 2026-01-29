-- 修改发货单表，添加物流和时间节点字段
ALTER TABLE shipment
    -- 目的地仓库关联
    ADD COLUMN destination_warehouse_id BIGINT UNSIGNED DEFAULT NULL COMMENT '目的地仓库ID' AFTER warehouse_id,
    ADD INDEX idx_destination_warehouse (destination_warehouse_id),

    -- 物流供应商关联
    ADD COLUMN logistics_provider_id BIGINT UNSIGNED DEFAULT NULL COMMENT '物流供应商ID' AFTER destination_address,
    ADD INDEX idx_logistics_provider (logistics_provider_id),

    -- 运费报价关联
    ADD COLUMN shipping_rate_id BIGINT UNSIGNED DEFAULT NULL COMMENT '运费报价ID' AFTER logistics_provider_id,
    ADD INDEX idx_shipping_rate (shipping_rate_id),

    -- 运输方式
    ADD COLUMN transport_mode ENUM('EXPRESS', 'AIR', 'SEA', 'RAIL', 'TRUCK') DEFAULT NULL COMMENT '运输方式' AFTER shipping_rate_id,
    ADD INDEX idx_transport_mode (transport_mode),

    -- 详细时间节点（带时分秒）
    ADD COLUMN confirmed_at DATETIME DEFAULT NULL COMMENT '确认时间' AFTER status,
    ADD COLUMN shipped_at DATETIME DEFAULT NULL COMMENT '发货时间' AFTER confirmed_at,
    ADD COLUMN delivered_at DATETIME DEFAULT NULL COMMENT '送达时间' AFTER shipped_at,

    -- 操作人记录
    ADD COLUMN confirmed_by BIGINT UNSIGNED DEFAULT NULL COMMENT '确认人' AFTER delivered_at,
    ADD COLUMN shipped_by BIGINT UNSIGNED DEFAULT NULL COMMENT '发货人' AFTER confirmed_by,
    ADD COLUMN delivered_by BIGINT UNSIGNED DEFAULT NULL COMMENT '签收人' AFTER shipped_by,

    -- 外键约束
    ADD CONSTRAINT fk_shipment_destination_warehouse
        FOREIGN KEY (destination_warehouse_id) REFERENCES warehouse(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_shipment_logistics_provider
        FOREIGN KEY (logistics_provider_id) REFERENCES logistics_provider(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_shipment_shipping_rate
        FOREIGN KEY (shipping_rate_id) REFERENCES shipping_rate(id) ON DELETE RESTRICT;

-- 创建索引用于查询时间范围
ALTER TABLE shipment
    ADD INDEX idx_confirmed_at (confirmed_at),
    ADD INDEX idx_shipped_at (shipped_at),
    ADD INDEX idx_delivered_at (delivered_at);
