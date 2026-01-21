-- ============================================
-- 产品图片管理 i18n 文案
-- ============================================
INSERT INTO field_label (label_key, module, scene, status, labels) VALUES
('product.images.title', 'product', 'images', 'active', JSON_OBJECT('zh-CN', '图片管理', 'en-US', 'Image Manager')),
('product.images.save', 'product', 'images', 'active', JSON_OBJECT('zh-CN', '提交保存', 'en-US', 'Save')),
('product.images.back', 'product', 'images', 'active', JSON_OBJECT('zh-CN', '返回列表', 'en-US', 'Back to List')),
('product.images.upload', 'product', 'images', 'active', JSON_OBJECT('zh-CN', '上传图片', 'en-US', 'Upload Images')),
('product.images.batchUpload', 'product', 'images', 'active', JSON_OBJECT('zh-CN', '批量上传', 'en-US', 'Batch Upload')),
('product.images.singleUpload', 'product', 'images', 'active', JSON_OBJECT('zh-CN', '单张上传', 'en-US', 'Upload')),
('product.images.limit', 'product', 'images', 'active', JSON_OBJECT('zh-CN', '最多10张，单张不超过10MB（jpg/jpeg/png/webp）', 'en-US', 'Up to 10 images, max 10MB each (jpg/jpeg/png/webp)')),
('product.images.dragHint', 'product', 'images', 'active', JSON_OBJECT('zh-CN', '拖拽图片调整顺序，第一张为主图', 'en-US', 'Drag to reorder, first image is primary')),
('product.images.count', 'product', 'images', 'active', JSON_OBJECT('zh-CN', '已上传 {count}/10', 'en-US', 'Uploaded {count}/10')),
('product.images.primary', 'product', 'images', 'active', JSON_OBJECT('zh-CN', '主图', 'en-US', 'Primary')),
('product.images.remove', 'product', 'images', 'active', JSON_OBJECT('zh-CN', '删除', 'en-US', 'Remove')),
('product.images.dirty', 'product', 'images', 'active', JSON_OBJECT('zh-CN', '未保存', 'en-US', 'Unsaved')),
('product.images.leaveConfirm', 'product', 'images', 'active', JSON_OBJECT('zh-CN', '有未保存的图片更改，确认离开？', 'en-US', 'You have unsaved changes. Leave page?')),
('product.images.warning', 'product', 'images', 'active', JSON_OBJECT('zh-CN', '提示', 'en-US', 'Warning'))
ON DUPLICATE KEY UPDATE
  module = VALUES(module),
  scene = VALUES(scene),
  status = VALUES(status),
  labels = VALUES(labels);
