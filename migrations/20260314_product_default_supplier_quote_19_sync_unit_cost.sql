-- 先用历史产品成本补齐“默认供应商基础报价”中的空价格，避免把已有有效成本回刷成 0
UPDATE product_supplier_quote q
JOIN product p
  ON p.id = q.product_id
 AND p.supplier_id = q.supplier_id
SET q.price = p.unit_cost
WHERE p.supplier_id IS NOT NULL
  AND p.unit_cost IS NOT NULL
  AND p.unit_cost > 0
  AND q.price <= 0;

-- 再把产品缓存成本统一回刷为默认供应商报价，避免 product.unit_cost 与默认报价长期分叉
UPDATE product p
JOIN product_supplier_quote q
  ON q.product_id = p.id
 AND q.supplier_id = p.supplier_id
SET p.unit_cost = q.price
WHERE p.supplier_id IS NOT NULL
  AND (p.unit_cost IS NULL OR p.unit_cost <> q.price);
