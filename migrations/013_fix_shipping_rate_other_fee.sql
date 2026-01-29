-- 修正运费报价表的其他费用字段
-- 将 other_charges (VARCHAR) 改为 other_fee (DECIMAL)

ALTER TABLE shipping_rate
    DROP COLUMN IF EXISTS other_charges;

ALTER TABLE shipping_rate
    ADD COLUMN other_fee DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '其他费用' AFTER currency;
