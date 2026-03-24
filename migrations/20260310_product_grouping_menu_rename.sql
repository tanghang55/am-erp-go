UPDATE `menu`
SET
  `title` = '产品归组',
  `title_en` = 'Product Grouping',
  `path` = '/product/groups',
  `gmt_modified` = NOW()
WHERE `code` = 'PRODUCT_PARENTS';
