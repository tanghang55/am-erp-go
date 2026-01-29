-- 创建物流服务表
CREATE TABLE IF NOT EXISTS logistics_service (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    service_code VARCHAR(50) NOT NULL UNIQUE COMMENT '服务代码（如：SLOW_SHIP、FAST_SHIP）',
    service_name VARCHAR(100) NOT NULL COMMENT '服务名称（如：慢船、快船、美森快船）',
    transport_mode ENUM('EXPRESS','AIR','SEA','RAIL','TRUCK') NOT NULL COMMENT '运输方式',
    destination_region VARCHAR(100) DEFAULT NULL COMMENT '目的地站点/国家（如：美国、欧洲、日本）',
    description TEXT DEFAULT NULL COMMENT '服务描述',
    status ENUM('ACTIVE','INACTIVE') NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
    created_by BIGINT UNSIGNED DEFAULT NULL COMMENT '创建人ID',
    updated_by BIGINT UNSIGNED DEFAULT NULL COMMENT '更新人ID',
    gmt_create TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    gmt_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    INDEX idx_transport_mode (transport_mode),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='物流服务表';

-- 修改运费报价表，添加service_id外键
ALTER TABLE shipping_rate
    ADD COLUMN service_id BIGINT UNSIGNED DEFAULT NULL COMMENT '物流服务ID' AFTER transport_mode,
    ADD INDEX idx_service_id (service_id);

-- 插入一些示例数据
INSERT INTO logistics_service (service_code, service_name, transport_mode, destination_region, description, status) VALUES
('SEA_SLOW_US', '美国慢船', 'SEA', '美国', '经济型海运服务，时效较慢但价格优惠', 'ACTIVE'),
('SEA_FAST_US', '美国快船', 'SEA', '美国', '快速海运服务，时效较快', 'ACTIVE'),
('SEA_MATSON_US', '美森快船', 'SEA', '美国', 'Matson快船服务，时效最快', 'ACTIVE'),
('SEA_SLOW_EU', '欧洲慢船', 'SEA', '欧洲', '经济型海运服务到欧洲', 'ACTIVE'),
('AIR_STANDARD', '标准空运', 'AIR', '全球', '标准空运服务', 'ACTIVE'),
('AIR_EXPRESS', '快速空运', 'AIR', '全球', '加急空运服务', 'ACTIVE'),
('EXPRESS_DHL', 'DHL快递', 'EXPRESS', '全球', 'DHL国际快递服务', 'ACTIVE'),
('EXPRESS_FEDEX', 'FedEx快递', 'EXPRESS', '全球', 'FedEx国际快递服务', 'ACTIVE');
