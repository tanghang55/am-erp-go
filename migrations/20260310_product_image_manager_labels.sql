INSERT INTO `field_label` (`label_key`, `module`, `scene`, `status`, `remark`, `labels`)
VALUES
  ('product.images.boardTitle', 'product', 'images', 'active', NULL, JSON_OBJECT('en-US', 'Image Board', 'zh-CN', '图片编排')),
  ('product.images.boardSubtitle', 'product', 'images', 'active', NULL, JSON_OBJECT('en-US', 'Drag to reorder. The first image is used as the primary image.', 'zh-CN', '拖拽排序，首图会作为产品主图展示。')),
  ('product.images.metricTotal', 'product', 'images', 'active', NULL, JSON_OBJECT('en-US', 'Uploaded', 'zh-CN', '已上传')),
  ('product.images.metricRemaining', 'product', 'images', 'active', NULL, JSON_OBJECT('en-US', 'Remaining Slots', 'zh-CN', '剩余槽位')),
  ('product.images.metricPrimary', 'product', 'images', 'active', NULL, JSON_OBJECT('en-US', 'Primary Image', 'zh-CN', '当前主图')),
  ('product.images.metricTotalHint', 'product', 'images', 'active', NULL, JSON_OBJECT('en-US', 'Current occupied image slots', 'zh-CN', '当前已占用图片位')),
  ('product.images.metricRemainingHint', 'product', 'images', 'active', NULL, JSON_OBJECT('en-US', 'Available upload capacity', 'zh-CN', '还可以继续上传')),
  ('product.images.metricPrimaryHint', 'product', 'images', 'active', NULL, JSON_OBJECT('en-US', 'The first image syncs as the product primary', 'zh-CN', '首图会同步到产品主图')),
  ('product.images.primaryReady', 'product', 'images', 'active', NULL, JSON_OBJECT('en-US', 'Configured', 'zh-CN', '已设置')),
  ('product.images.primaryEmpty', 'product', 'images', 'active', NULL, JSON_OBJECT('en-US', 'Not Set', 'zh-CN', '未设置')),
  ('product.images.primaryPanelTitle', 'product', 'images', 'active', NULL, JSON_OBJECT('en-US', 'Primary Preview', 'zh-CN', '主图预览')),
  ('product.images.primaryPanelSubtitle', 'product', 'images', 'active', NULL, JSON_OBJECT('en-US', 'Use this panel to verify the primary image and image status.', 'zh-CN', '右侧用于快速核对首图和图片状态。')),
  ('product.images.ruleTitle', 'product', 'images', 'active', NULL, JSON_OBJECT('en-US', 'Upload Rules', 'zh-CN', '上传规则')),
  ('product.images.ruleFill', 'product', 'images', 'active', NULL, JSON_OBJECT('en-US', 'Fill the front slots first to keep the primary order clear.', 'zh-CN', '建议优先补齐前排图片槽位，避免主图顺序混乱。'))
ON DUPLICATE KEY UPDATE
  `module` = VALUES(`module`),
  `scene` = VALUES(`scene`),
  `status` = VALUES(`status`),
  `remark` = VALUES(`remark`),
  `labels` = VALUES(`labels`),
  `gmt_modified` = CURRENT_TIMESTAMP;
