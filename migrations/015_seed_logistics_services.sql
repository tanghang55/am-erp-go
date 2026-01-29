-- 添加测试物流服务数据

INSERT INTO logistics_service (
    service_code,
    service_name,
    transport_mode,
    destination_region,
    description,
    status
) VALUES
-- 快递服务
('EXPRESS-DOMESTIC', '国内快递', 'EXPRESS', '中国大陆', '国内标准快递服务，1-3天送达', 'ACTIVE'),
('EXPRESS-INTL-STD', '国际标准快递', 'EXPRESS', '全球', '国际标准快递，3-7个工作日', 'ACTIVE'),
('EXPRESS-INTL-PRIORITY', '国际优先快递', 'EXPRESS', '全球', '国际优先快递，2-4个工作日', 'ACTIVE'),
('EXPRESS-USA', '美国快递专线', 'EXPRESS', '美国', '美国快递专线，3-5个工作日', 'ACTIVE'),
('EXPRESS-EU', '欧洲快递专线', 'EXPRESS', '欧洲', '欧洲快递专线，4-7个工作日', 'ACTIVE'),

-- 空运服务
('AIR-STANDARD', '标准空运', 'AIR', '全球', '标准空运服务，5-10个工作日', 'ACTIVE'),
('AIR-FAST', '快速空运', 'AIR', '全球', '快速空运服务，3-5个工作日', 'ACTIVE'),
('AIR-USA', '美国空运专线', 'AIR', '美国', '美国空运专线，4-6个工作日', 'ACTIVE'),
('AIR-EU', '欧洲空运专线', 'AIR', '欧洲', '欧洲空运专线，5-8个工作日', 'ACTIVE'),

-- 海运服务
('SEA-FCL', '整柜海运', 'SEA', '全球', '整柜海运服务，20-40天', 'ACTIVE'),
('SEA-LCL', '拼箱海运', 'SEA', '全球', '拼箱海运服务，25-45天', 'ACTIVE'),
('SEA-USA-WEST', '美西海运', 'SEA', '美国西海岸', '美国西海岸海运，15-20天', 'ACTIVE'),
('SEA-USA-EAST', '美东海运', 'SEA', '美国东海岸', '美国东海岸海运，25-30天', 'ACTIVE'),
('SEA-EU', '欧洲海运', 'SEA', '欧洲', '欧洲海运服务，30-40天', 'ACTIVE'),

-- 铁路运输
('RAIL-EU', '中欧班列', 'RAIL', '欧洲', '中欧铁路班列，15-20天', 'ACTIVE'),
('RAIL-RUSSIA', '中俄班列', 'RAIL', '俄罗斯', '中俄铁路班列，10-15天', 'ACTIVE'),

-- 卡车运输
('TRUCK-DOMESTIC', '国内卡车运输', 'TRUCK', '中国大陆', '国内卡车整车/零担运输', 'ACTIVE'),
('TRUCK-CROSS-BORDER', '跨境卡车运输', 'TRUCK', '东南亚', '东南亚跨境卡车运输', 'ACTIVE');
