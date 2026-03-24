INSERT INTO `menu` (
  `title`, `title_en`, `code`, `parent_id`, `path`, `icon`, `component`, `sort`,
  `is_hidden`, `permission_code`, `status`, `gmt_create`, `gmt_modified`
)
SELECT
  '产品配置', 'Product Config', 'PRODUCT_CONFIG',
  `id`, '/product/config', NULL, NULL, 20,
  0, 'product.manage', 'ACTIVE', NOW(), NOW()
FROM `menu`
WHERE `code` = 'PRODUCT'
LIMIT 1;
