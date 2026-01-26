-- 006: 装箱规格添加每箱数量字段
-- 用于计算发货产品数量 = 箱数 × 每箱数量

ALTER TABLE package_spec ADD COLUMN quantity_per_box INT UNSIGNED NOT NULL DEFAULT 1 COMMENT '每箱产品数量' AFTER weight;
