-- 添加其他费用字段到运费报价表
ALTER TABLE shipping_rate
    ADD COLUMN other_charges VARCHAR(500) DEFAULT NULL COMMENT '其他费用说明（如：燃油附加费15%、偏远地区+50元）' AFTER currency;
