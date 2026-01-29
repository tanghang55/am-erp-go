-- 物流供应商表
CREATE TABLE IF NOT EXISTS logistics_provider (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',

    -- 基本信息
    provider_code VARCHAR(50) NOT NULL COMMENT '供应商代码',
    provider_name VARCHAR(200) NOT NULL COMMENT '供应商名称',
    provider_type ENUM('FREIGHT_FORWARDER', 'COURIER', 'SHIPPING_LINE', 'AIRLINE') NOT NULL COMMENT '供应商类型: 货代/快递/船公司/航空',

    -- 服务类型
    service_types VARCHAR(200) DEFAULT NULL COMMENT '服务类型(EXPRESS,AIR,SEA,RAIL),逗号分隔',

    -- 联系信息
    contact_person VARCHAR(100) DEFAULT NULL COMMENT '联系人',
    contact_phone VARCHAR(50) DEFAULT NULL COMMENT '联系电话',
    contact_email VARCHAR(100) DEFAULT NULL COMMENT '联系邮箱',
    website VARCHAR(200) DEFAULT NULL COMMENT '网站',

    -- 地址
    country VARCHAR(50) DEFAULT NULL COMMENT '国家',
    city VARCHAR(100) DEFAULT NULL COMMENT '城市',
    address TEXT DEFAULT NULL COMMENT '地址',

    -- 账号信息
    account_number VARCHAR(100) DEFAULT NULL COMMENT '客户账号',
    credit_days INT DEFAULT 0 COMMENT '账期天数',

    -- 状态
    status ENUM('ACTIVE', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',

    -- 备注
    remark TEXT DEFAULT NULL COMMENT '备注',

    -- 审计字段
    created_by BIGINT UNSIGNED DEFAULT NULL COMMENT '创建人',
    updated_by BIGINT UNSIGNED DEFAULT NULL COMMENT '更新人',
    gmt_create DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',

    UNIQUE KEY uk_provider_code (provider_code),
    INDEX idx_provider_type (provider_type),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='物流供应商表';
