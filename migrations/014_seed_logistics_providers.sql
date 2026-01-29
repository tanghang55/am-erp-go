-- 添加测试物流供应商数据

INSERT INTO logistics_provider (
    provider_code,
    provider_name,
    provider_type,
    service_types,
    contact_person,
    contact_phone,
    contact_email,
    country,
    city,
    status
) VALUES
-- 快递公司
('DHL-CN', 'DHL中国', 'COURIER', 'EXPRESS,AIR', '张三', '400-810-8000', 'service@dhl.com.cn', '中国', '上海', 'ACTIVE'),
('FEDEX-CN', 'FedEx中国', 'COURIER', 'EXPRESS,AIR', '李四', '400-886-1888', 'service@fedex.com.cn', '中国', '上海', 'ACTIVE'),
('UPS-CN', 'UPS中国', 'COURIER', 'EXPRESS,AIR', '王五', '400-820-8388', 'service@ups.com.cn', '中国', '上海', 'ACTIVE'),
('SF-EXPRESS', '顺丰速运', 'COURIER', 'EXPRESS,AIR', '赵六', '95338', 'service@sf-express.com', '中国', '深圳', 'ACTIVE'),

-- 货代公司
('FREIGHT-FWD-01', '环球货运代理', 'FREIGHT_FORWARDER', 'AIR,SEA', '钱七', '021-12345678', 'info@global-freight.com', '中国', '上海', 'ACTIVE'),
('FREIGHT-FWD-02', '中远海运物流', 'FREIGHT_FORWARDER', 'SEA,RAIL', '孙八', '0755-88888888', 'service@cosco-logistics.com', '中国', '深圳', 'ACTIVE'),
('FREIGHT-FWD-03', '嘉里物流', 'FREIGHT_FORWARDER', 'AIR,SEA,RAIL,TRUCK', '周九', '400-820-3333', 'service@kerrylogistics.com', '中国', '上海', 'ACTIVE'),

-- 船公司
('MAERSK', '马士基航运', 'SHIPPING_LINE', 'SEA', '吴十', '400-120-6888', 'service@maersk.com.cn', '中国', '上海', 'ACTIVE'),
('MSC', '地中海航运', 'SHIPPING_LINE', 'SEA', '郑十一', '021-63303000', 'service@msc.com.cn', '中国', '上海', 'ACTIVE'),
('COSCO', '中远海运集运', 'SHIPPING_LINE', 'SEA', '林十二', '400-820-8888', 'service@coscoshipping.com', '中国', '上海', 'ACTIVE'),

-- 航空公司
('CA-CARGO', '国航货运', 'AIRLINE', 'AIR', '陈十三', '010-95583', 'cargo@airchina.com', '中国', '北京', 'ACTIVE'),
('CZ-CARGO', '南航货运', 'AIRLINE', 'AIR', '黄十四', '020-95539', 'cargo@csair.com', '中国', '广州', 'ACTIVE');
