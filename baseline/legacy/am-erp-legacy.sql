/*
 Navicat Premium Dump SQL

 Source Server         : localhost_3306
 Source Server Type    : MySQL
 Source Server Version : 80044 (8.0.44)
 Source Host           : localhost:3306
 Source Schema         : am-erp

 Target Server Type    : MySQL
 Target Server Version : 80044 (8.0.44)
 File Encoding         : 65001

 Date: 29/01/2026 17:01:19
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for audit_log
-- ----------------------------
DROP TABLE IF EXISTS `audit_log`;
CREATE TABLE `audit_log`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `trace_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '追踪ID（全链路追踪）',
  `user_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '操作人ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作人用户名（冗余，防止用户删除后无法追溯）',
  `module` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模块名（Catalog/Inventory/Procurement/Shipping/Packaging/Finance/Imports）',
  `action` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作动作（CREATE/UPDATE/DELETE/IMPORT/BIND/RECEIVE等）',
  `entity_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '实体类型（SKU/Listing/PO/Shipment/Movement等）',
  `entity_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '实体ID',
  `changes` json NULL COMMENT '变更内容（{\"before\": {...}, \"after\": {...}}）',
  `ip_address` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'IP地址（支持IPv4和IPv6）',
  `user_agent` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT 'User Agent',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_trace_id`(`trace_id` ASC) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_module`(`module` ASC) USING BTREE,
  INDEX `idx_action`(`action` ASC) USING BTREE,
  INDEX `idx_entity_type_entity_id`(`entity_type` ASC, `entity_id` ASC) USING BTREE,
  INDEX `idx_gmt_create`(`gmt_create` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 123 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '审计日志表（记录谁在何时对什么做了什么操作）' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of audit_log
-- ----------------------------
INSERT INTO `audit_log` VALUES (1, '29bb1540-f464-11f0-92b5-8c32235251f8', NULL, 'SYSTEM', 'Identity', 'INIT', 'User', '1', '{\"message\": \"System initialized with default admin user\"}', NULL, NULL, '2026-01-18 19:52:28', '2026-01-18 19:52:28');
INSERT INTO `audit_log` VALUES (2, '8727d524-f464-11f0-92b5-8c32235251f8', NULL, 'SYSTEM', 'Identity', 'INIT', 'User', '1', '{\"message\": \"System initialized with default admin user\"}', NULL, NULL, '2026-01-18 19:55:04', '2026-01-18 19:55:04');
INSERT INTO `audit_log` VALUES (3, 'f2e7eaf1-f464-11f0-92b5-8c32235251f8', NULL, 'SYSTEM', 'Identity', 'INIT', 'User', '1', '{\"message\": \"System initialized with default admin user\"}', NULL, NULL, '2026-01-18 19:58:05', '2026-01-18 19:58:05');
INSERT INTO `audit_log` VALUES (4, '919926b3-6f5b-4813-a4d7-71b04a7433a8', 8, 'admin', 'Product', 'CREATE', 'ProductSupplierQuote', '1', '{\"after\": {\"id\": 1, \"price\": 1, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-20T18:27:31.051+08:00\", \"product_id\": 6, \"supplier_id\": 2, \"gmt_modified\": \"2026-01-20T18:27:31.051+08:00\", \"lead_time_days\": 0}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-20 18:27:31', '2026-01-20 18:27:31');
INSERT INTO `audit_log` VALUES (5, 'fd0e7d17-c403-4d4b-a6ae-e2dc46e2162b', 8, 'admin', 'Product', 'CREATE', 'ProductSupplierQuote', '2', '{\"after\": {\"id\": 2, \"price\": 3, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-20T18:27:42.846+08:00\", \"product_id\": 6, \"supplier_id\": 3, \"gmt_modified\": \"2026-01-20T18:27:42.846+08:00\", \"lead_time_days\": 0}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-20 18:27:43', '2026-01-20 18:27:43');
INSERT INTO `audit_log` VALUES (6, 'bff6e4b6-0ed2-4b62-9930-05d21c2e1379', 8, 'admin', 'Product', 'UPDATE', 'Product', '6', '{\"after\": {\"default_supplier_id\": 2}, \"before\": {\"default_supplier_id\": 1}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-20 18:35:33', '2026-01-20 18:35:33');
INSERT INTO `audit_log` VALUES (7, '17358b8f-9fdf-47fd-8337-3ea2ec352653', 8, 'admin', 'Product', 'UPDATE', 'Product', '6', '{\"after\": {\"default_supplier_id\": 3}, \"before\": {\"default_supplier_id\": 2}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-20 18:35:41', '2026-01-20 18:35:41');
INSERT INTO `audit_log` VALUES (8, '2363fdde-e68f-473f-9ec7-46b289e5f9ce', 8, 'admin', 'Product', 'CREATE', 'ProductSupplierQuote', '3', '{\"after\": {\"id\": 3, \"price\": 0, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T15:09:41.243+08:00\", \"product_id\": 1, \"supplier_id\": 1, \"gmt_modified\": \"2026-01-21T15:09:41.243+08:00\", \"lead_time_days\": 0}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 15:09:41', '2026-01-21 15:09:41');
INSERT INTO `audit_log` VALUES (9, 'd9ca2a56-7978-47e4-a9f8-ba484ccf8960', 8, 'admin', 'Product', 'CREATE', 'ProductSupplierQuote', '4', '{\"after\": {\"id\": 4, \"price\": 0, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T15:09:45.675+08:00\", \"product_id\": 2, \"supplier_id\": 1, \"gmt_modified\": \"2026-01-21T15:09:45.675+08:00\", \"lead_time_days\": 0}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 15:09:46', '2026-01-21 15:09:46');
INSERT INTO `audit_log` VALUES (10, 'e3fec727-f37c-44fb-a00c-b3661d175dbc', 8, 'admin', 'Product', 'CREATE', 'ProductSupplierQuote', '5', '{\"after\": {\"id\": 5, \"price\": 0, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T15:09:49.3+08:00\", \"product_id\": 3, \"supplier_id\": 2, \"gmt_modified\": \"2026-01-21T15:09:49.3+08:00\", \"lead_time_days\": 0}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 15:09:49', '2026-01-21 15:09:49');
INSERT INTO `audit_log` VALUES (11, '8f3ef0e6-edaf-4efc-86d1-f9689ee3a396', 8, 'admin', 'Product', 'CREATE', 'ProductSupplierQuote', '6', '{\"after\": {\"id\": 6, \"price\": 0, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T15:09:53.138+08:00\", \"product_id\": 5, \"supplier_id\": 3, \"gmt_modified\": \"2026-01-21T15:09:53.138+08:00\", \"lead_time_days\": 0}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 15:09:53', '2026-01-21 15:09:53');
INSERT INTO `audit_log` VALUES (12, 'cf2bbcd4-9372-48b2-939c-6d2505d59edb', 8, 'admin', 'Product', 'CREATE', 'ProductSupplierQuote', '7', '{\"after\": {\"id\": 7, \"price\": 0, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T15:09:56.953+08:00\", \"product_id\": 4, \"supplier_id\": 4, \"gmt_modified\": \"2026-01-21T15:09:56.953+08:00\", \"lead_time_days\": 0}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 15:09:57', '2026-01-21 15:09:57');
INSERT INTO `audit_log` VALUES (13, '540c7a48-3505-499d-9f7e-24863594b854', 8, 'admin', 'Product', 'CREATE', 'ProductSupplierQuote', '8', '{\"after\": {\"id\": 8, \"price\": 0, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T18:21:24.203+08:00\", \"product_id\": 1, \"supplier_id\": 2, \"gmt_modified\": \"2026-01-21T18:21:24.203+08:00\", \"lead_time_days\": 0}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 18:21:24', '2026-01-21 18:21:24');
INSERT INTO `audit_log` VALUES (14, '8d3d306f-3f95-4630-89a1-999e0bb9dfce', 8, 'admin', 'Product', 'CREATE', 'ProductSupplierQuote', '9', '{\"after\": {\"id\": 9, \"price\": 0, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T18:21:29.689+08:00\", \"product_id\": 1, \"supplier_id\": 4, \"gmt_modified\": \"2026-01-21T18:21:29.689+08:00\", \"lead_time_days\": 0}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 18:21:30', '2026-01-21 18:21:30');
INSERT INTO `audit_log` VALUES (15, '9acdcf11-f117-45c1-9c95-01f8618f0262', 8, 'admin', 'Product', 'UPDATE', 'ProductSupplierQuote', '9', '{\"after\": {\"id\": 9, \"price\": 1, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T18:21:30+08:00\", \"product_id\": 1, \"supplier_id\": 4, \"gmt_modified\": \"2026-01-21T18:30:11.532+08:00\", \"lead_time_days\": 0}, \"before\": {\"id\": 9, \"price\": 0, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T18:21:30+08:00\", \"product_id\": 1, \"supplier_id\": 4, \"gmt_modified\": \"2026-01-21T18:21:30+08:00\", \"lead_time_days\": 0}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 18:30:12', '2026-01-21 18:30:12');
INSERT INTO `audit_log` VALUES (16, '8b4090fa-bc17-4959-8764-0796900839b1', 8, 'admin', 'Product', 'UPDATE', 'ProductSupplierQuote', '8', '{\"after\": {\"id\": 8, \"price\": 10, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 10, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T18:21:24+08:00\", \"product_id\": 1, \"supplier_id\": 2, \"gmt_modified\": \"2026-01-21T18:30:21.585+08:00\", \"lead_time_days\": 0}, \"before\": {\"id\": 8, \"price\": 0, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T18:21:24+08:00\", \"product_id\": 1, \"supplier_id\": 2, \"gmt_modified\": \"2026-01-21T18:21:24+08:00\", \"lead_time_days\": 0}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 18:30:22', '2026-01-21 18:30:22');
INSERT INTO `audit_log` VALUES (17, 'cc5c15a1-3b39-4f5e-bb6e-cd731623d163', 8, 'admin', 'Product', 'UPDATE', 'ProductSupplierQuote', '3', '{\"after\": {\"id\": 3, \"price\": 100, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 100, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T15:09:41+08:00\", \"product_id\": 1, \"supplier_id\": 1, \"gmt_modified\": \"2026-01-21T18:30:33.664+08:00\", \"lead_time_days\": 0}, \"before\": {\"id\": 3, \"price\": 0, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T15:09:41+08:00\", \"product_id\": 1, \"supplier_id\": 1, \"gmt_modified\": \"2026-01-21T15:09:41+08:00\", \"lead_time_days\": 0}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 18:30:34', '2026-01-21 18:30:34');
INSERT INTO `audit_log` VALUES (18, '54dce193-4294-4a2d-917d-f3399d8942ee', 8, 'admin', 'Product', 'UPDATE', 'ProductSupplierQuote', '4', '{\"after\": {\"id\": 4, \"price\": 1, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T15:09:46+08:00\", \"product_id\": 2, \"supplier_id\": 1, \"gmt_modified\": \"2026-01-21T18:37:41.342+08:00\", \"lead_time_days\": 0}, \"before\": {\"id\": 4, \"price\": 0, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T15:09:46+08:00\", \"product_id\": 2, \"supplier_id\": 1, \"gmt_modified\": \"2026-01-21T15:09:46+08:00\", \"lead_time_days\": 0}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 18:37:41', '2026-01-21 18:37:41');
INSERT INTO `audit_log` VALUES (19, 'da734b17-d09d-4523-ab0d-b504a4cea5a3', 8, 'admin', 'Product', 'UPDATE', 'ProductSupplierQuote', '5', '{\"after\": {\"id\": 5, \"price\": 1, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T15:09:49+08:00\", \"product_id\": 3, \"supplier_id\": 2, \"gmt_modified\": \"2026-01-21T18:37:47.417+08:00\", \"lead_time_days\": 0}, \"before\": {\"id\": 5, \"price\": 0, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T15:09:49+08:00\", \"product_id\": 3, \"supplier_id\": 2, \"gmt_modified\": \"2026-01-21T15:09:49+08:00\", \"lead_time_days\": 0}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 18:37:47', '2026-01-21 18:37:47');
INSERT INTO `audit_log` VALUES (20, '7028ae38-ccb2-43ec-ae78-ad5757bc8256', 8, 'admin', 'Product', 'UPDATE', 'ProductSupplierQuote', '6', '{\"after\": {\"id\": 6, \"price\": 3, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T15:09:53+08:00\", \"product_id\": 5, \"supplier_id\": 3, \"gmt_modified\": \"2026-01-21T18:37:52.61+08:00\", \"lead_time_days\": 0}, \"before\": {\"id\": 6, \"price\": 0, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T15:09:53+08:00\", \"product_id\": 5, \"supplier_id\": 3, \"gmt_modified\": \"2026-01-21T15:09:53+08:00\", \"lead_time_days\": 0}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 18:37:53', '2026-01-21 18:37:53');
INSERT INTO `audit_log` VALUES (21, 'e5d73758-770e-4677-8503-b2cdc861d1f0', 8, 'admin', 'Product', 'UPDATE', 'ProductSupplierQuote', '7', '{\"after\": {\"id\": 7, \"price\": 2, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T15:09:57+08:00\", \"product_id\": 4, \"supplier_id\": 4, \"gmt_modified\": \"2026-01-21T18:37:58.75+08:00\", \"lead_time_days\": 0}, \"before\": {\"id\": 7, \"price\": 0, \"remark\": \"\", \"status\": \"ACTIVE\", \"qty_moq\": 1, \"currency\": \"USD\", \"gmt_create\": \"2026-01-21T15:09:57+08:00\", \"product_id\": 4, \"supplier_id\": 4, \"gmt_modified\": \"2026-01-21T15:09:57+08:00\", \"lead_time_days\": 0}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 18:37:59', '2026-01-21 18:37:59');
INSERT INTO `audit_log` VALUES (22, '1a593c43-931c-4d88-8008-88947747bdf5', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '1', '{\"after\": {\"id\": 1, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601211838181700\", \"created_at\": \"2026-01-21T18:38:18.724+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T18:38:18.724+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 18:38:19', '2026-01-21 18:38:19');
INSERT INTO `audit_log` VALUES (23, '86cc7921-69f6-41d5-b380-2c2b80e48858', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '5', '{\"after\": {\"id\": 5, \"items\": [{\"id\": 0, \"sku_id\": 4, \"currency\": \"USD\", \"subtotal\": 2, \"unit_cost\": 2, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601211838286900\", \"created_at\": \"2026-01-21T18:38:28.928+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T18:38:28.928+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 4, \"total_amount\": 2}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 18:38:29', '2026-01-21 18:38:29');
INSERT INTO `audit_log` VALUES (24, '002380a8-d036-496a-a186-bc3f0c0cd942', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '9', '{\"after\": {\"id\": 9, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601211847321100\", \"created_at\": \"2026-01-21T18:47:32.311+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T18:47:32.311+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 18:47:32', '2026-01-21 18:47:32');
INSERT INTO `audit_log` VALUES (25, '8968c4da-a65a-4a6a-a6c3-4edf24d201c5', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '11', '{\"after\": {\"id\": 11, \"items\": [{\"id\": 0, \"sku_id\": 5, \"currency\": \"USD\", \"subtotal\": 3, \"unit_cost\": 3, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 6, \"currency\": \"USD\", \"subtotal\": 3, \"unit_cost\": 3, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601211847322400\", \"created_at\": \"2026-01-21T18:47:32.312+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T18:47:32.312+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 3, \"total_amount\": 6}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 18:47:32', '2026-01-21 18:47:32');
INSERT INTO `audit_log` VALUES (26, '215a2c83-4628-4a27-8066-998f97466b38', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '12', '{\"after\": {\"id\": 12, \"items\": [{\"id\": 0, \"sku_id\": 4, \"currency\": \"USD\", \"subtotal\": 2, \"unit_cost\": 2, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601211847323400\", \"created_at\": \"2026-01-21T18:47:32.312+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T18:47:32.312+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 4, \"total_amount\": 2}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 18:47:32', '2026-01-21 18:47:32');
INSERT INTO `audit_log` VALUES (27, 'e68eee1a-c151-43c7-b122-6e602083bcd6', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '13', '{\"after\": {\"id\": 13, \"items\": [{\"id\": 0, \"sku_id\": 5, \"currency\": \"USD\", \"subtotal\": 3, \"unit_cost\": 3, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 6, \"currency\": \"USD\", \"subtotal\": 3, \"unit_cost\": 3, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601211847530400\", \"created_at\": \"2026-01-21T18:47:53.918+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T18:47:53.918+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 3, \"total_amount\": 6}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 18:47:54', '2026-01-21 18:47:54');
INSERT INTO `audit_log` VALUES (28, '7917afc2-b2e2-4ab8-b096-be9080b870c1', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '17', '{\"after\": {\"id\": 17, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601211853227800\", \"created_at\": \"2026-01-21T18:53:22.074+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T18:53:22.074+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 18:53:22', '2026-01-21 18:53:22');
INSERT INTO `audit_log` VALUES (29, '36bbe8f6-f0a0-43c1-b700-35f5c9891eb8', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '19', '{\"after\": {\"id\": 19, \"items\": [{\"id\": 0, \"sku_id\": 4, \"currency\": \"USD\", \"subtotal\": 2, \"unit_cost\": 2, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601211853225000\", \"created_at\": \"2026-01-21T18:53:22.074+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T18:53:22.074+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 4, \"total_amount\": 2}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 18:53:22', '2026-01-21 18:53:22');
INSERT INTO `audit_log` VALUES (30, '3a808b2c-94e0-49d4-af37-cd9a7409bff4', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '22', '{\"after\": {\"id\": 22, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601211904568100\", \"created_at\": \"2026-01-21T19:04:56.856+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T19:04:56.856+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 19:04:57', '2026-01-21 19:04:57');
INSERT INTO `audit_log` VALUES (31, 'c9aa856d-cf68-4291-a4dc-1dbe44bc9fb4', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '21', '{\"after\": {\"id\": 21, \"items\": [{\"id\": 0, \"sku_id\": 4, \"currency\": \"USD\", \"subtotal\": 2, \"unit_cost\": 2, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601211904565500\", \"created_at\": \"2026-01-21T19:04:56.856+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T19:04:56.856+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 4, \"total_amount\": 2}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 19:04:57', '2026-01-21 19:04:57');
INSERT INTO `audit_log` VALUES (32, '2f49b144-f22b-4092-a79e-a70ec55e24a4', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '25', '{\"after\": {\"id\": 25, \"items\": [{\"id\": 0, \"sku_id\": 5, \"currency\": \"USD\", \"subtotal\": 3, \"unit_cost\": 3, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 6, \"currency\": \"USD\", \"subtotal\": 3, \"unit_cost\": 3, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601211905087700\", \"created_at\": \"2026-01-21T19:05:08.689+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T19:05:08.689+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 3, \"total_amount\": 6}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 19:05:09', '2026-01-21 19:05:09');
INSERT INTO `audit_log` VALUES (33, 'b860b208-fbc3-4978-922e-95d4de01a048', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '29', '{\"after\": {\"id\": 29, \"items\": [{\"id\": 0, \"sku_id\": 5, \"currency\": \"USD\", \"subtotal\": 3, \"unit_cost\": 3, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 6, \"currency\": \"USD\", \"subtotal\": 3, \"unit_cost\": 3, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601211905422000\", \"created_at\": \"2026-01-21T19:05:42.82+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T19:05:42.82+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 3, \"total_amount\": 6}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 19:05:43', '2026-01-21 19:05:43');
INSERT INTO `audit_log` VALUES (34, '2fdcf704-3851-4743-bc0d-5d632c070cbe', 8, 'admin', 'Procurement', 'SUBMIT', 'PurchaseOrder', '29', '{\"after\": {\"id\": 29, \"remark\": \"\", \"status\": \"ORDERED\", \"currency\": \"USD\", \"supplier\": {\"id\": 3, \"name\": \"上海物流服务\"}, \"po_number\": \"PO202601211905422000\", \"created_at\": \"2026-01-21T19:05:43+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-21T19:07:42.8402648+08:00\", \"shipped_at\": null, \"updated_at\": \"2026-01-21T19:07:42.84+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 3, \"total_amount\": 6}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 19:07:43', '2026-01-21 19:07:43');
INSERT INTO `audit_log` VALUES (35, 'b7c26c1c-b5d0-4b8b-8b11-2134c04a4d09', 8, 'admin', 'Procurement', 'SHIP', 'PurchaseOrder', '29', '{\"after\": {\"id\": 29, \"remark\": \"\", \"status\": \"SHIPPED\", \"currency\": \"USD\", \"supplier\": {\"id\": 3, \"name\": \"上海物流服务\"}, \"po_number\": \"PO202601211905422000\", \"created_at\": \"2026-01-21T19:05:43+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-21T19:07:43+08:00\", \"shipped_at\": \"2026-01-21T19:07:50.4754378+08:00\", \"updated_at\": \"2026-01-21T19:07:50.475+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 3, \"total_amount\": 6}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 19:07:50', '2026-01-21 19:07:50');
INSERT INTO `audit_log` VALUES (36, 'ccc9d6ed-1b9d-4e60-9285-9660ec966bdf', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '33', '{\"after\": {\"id\": 33, \"items\": [{\"id\": 0, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601211955082500\", \"created_at\": \"2026-01-21T19:55:08.711+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T19:55:08.711+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 19:55:09', '2026-01-21 19:55:09');
INSERT INTO `audit_log` VALUES (37, '022a05c7-c963-4f89-98fd-de7864bdd413', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '35', '{\"after\": {\"id\": 35, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601211955149100\", \"created_at\": \"2026-01-21T19:55:14.178+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T19:55:14.178+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 19:55:14', '2026-01-21 19:55:14');
INSERT INTO `audit_log` VALUES (38, 'f47d3eef-7a7a-44d3-8b43-d158fbfd43f1', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '37', '{\"after\": {\"id\": 37, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 199.9, \"unit_cost\": 19.99, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 10, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 99.9, \"unit_cost\": 9.99, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 10, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 129.9, \"unit_cost\": 12.99, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 10, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"���Բɹ���\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601211957070800\", \"created_at\": \"2026-01-21T19:57:07.83+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T19:57:07.83+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 4, \"total_amount\": 429.7}, \"before\": null}', '::1', 'curl/8.14.1', '2026-01-21 19:57:08', '2026-01-21 19:57:08');
INSERT INTO `audit_log` VALUES (39, '9be1e643-29d0-4915-9a6b-51ffd0710e3a', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '38', '{\"after\": {\"id\": 38, \"items\": [{\"id\": 0, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601211959315700\", \"created_at\": \"2026-01-21T19:59:31.121+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T19:59:31.121+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 19:59:31', '2026-01-21 19:59:31');
INSERT INTO `audit_log` VALUES (40, '13b2f1a8-5b2f-47dc-97af-e55e54aa50d5', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '40', '{\"after\": {\"id\": 40, \"items\": [{\"id\": 0, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601212000498500\", \"created_at\": \"2026-01-21T20:00:49.82+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T20:00:49.82+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 20:00:50', '2026-01-21 20:00:50');
INSERT INTO `audit_log` VALUES (41, 'b5802298-5fea-4145-9b59-2abc2feb5a61', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '41', '{\"after\": {\"id\": 41, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601212000499100\", \"created_at\": \"2026-01-21T20:00:49.827+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T20:00:49.827+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 20:00:50', '2026-01-21 20:00:50');
INSERT INTO `audit_log` VALUES (42, '0c1c7c1c-8b2a-4ab0-ae0c-014713d2e15e', 8, 'admin', 'Procurement', 'SUBMIT', 'PurchaseOrder', '40', '{\"after\": {\"id\": 40, \"remark\": \"\", \"status\": \"ORDERED\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO202601212000498500\", \"created_at\": \"2026-01-21T20:00:50+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-21T20:01:32.8074011+08:00\", \"shipped_at\": null, \"updated_at\": \"2026-01-21T20:01:32.807+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 20:01:33', '2026-01-21 20:01:33');
INSERT INTO `audit_log` VALUES (43, '17327091-7903-42ea-a2e2-480ab119e07a', 8, 'admin', 'Procurement', 'SHIP', 'PurchaseOrder', '40', '{\"after\": {\"id\": 40, \"remark\": \"\", \"status\": \"SHIPPED\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO202601212000498500\", \"created_at\": \"2026-01-21T20:00:50+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-21T20:01:33+08:00\", \"shipped_at\": \"2026-01-21T20:01:35.5802567+08:00\", \"updated_at\": \"2026-01-21T20:01:35.581+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 20:01:36', '2026-01-21 20:01:36');
INSERT INTO `audit_log` VALUES (44, '9743a53e-6efc-4c12-b84d-7d7e0ff5659f', 8, 'admin', 'Procurement', 'RECEIVE', 'PurchaseOrder', '40', '{\"after\": {\"id\": 40, \"items\": [{\"id\": 28, \"sku\": {\"id\": 3, \"title\": \"组合子产品B\", \"seller_sku\": \"SKU-CHILD-003\"}, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-21T20:00:50+08:00\", \"updated_at\": \"2026-01-21T20:00:50+08:00\", \"qty_ordered\": 1, \"qty_received\": 1, \"purchase_order_id\": 40}], \"remark\": \"\", \"status\": \"RECEIVED\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO202601212000498500\", \"created_at\": \"2026-01-21T20:00:50+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-21T20:01:33+08:00\", \"shipped_at\": \"2026-01-21T20:01:36+08:00\", \"updated_at\": \"2026-01-21T20:04:12.618+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": \"2026-01-21T20:04:12.6102033+08:00\", \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-21 20:04:13', '2026-01-21 20:04:13');
INSERT INTO `audit_log` VALUES (45, '1dd8754b-5e0d-4b6a-8bfe-747cd7585e4c', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '42', '{\"after\": {\"id\": 42, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601221716361900\", \"created_at\": \"2026-01-22T17:16:36.073+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-22T17:16:36.073+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-22 17:16:36', '2026-01-22 17:16:36');
INSERT INTO `audit_log` VALUES (46, 'd9f22e6f-ed1b-47fe-918d-b24bf48cfbe0', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '44', '{\"after\": {\"id\": 44, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601221716443200\", \"created_at\": \"2026-01-22T17:16:44.579+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-22T17:16:44.579+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-22 17:16:45', '2026-01-22 17:16:45');
INSERT INTO `audit_log` VALUES (47, '73b0abd5-894c-4552-b1dd-d0afa78f430d', 8, 'admin', 'Procurement', 'SUBMIT', 'PurchaseOrder', '44', '{\"after\": {\"id\": 44, \"remark\": \"\", \"status\": \"ORDERED\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601221716443200\", \"created_at\": \"2026-01-22T17:16:45+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-22T17:17:11.9149208+08:00\", \"shipped_at\": null, \"updated_at\": \"2026-01-22T17:17:11.915+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-22 17:17:12', '2026-01-22 17:17:12');
INSERT INTO `audit_log` VALUES (48, '5c1b2bb7-61ae-4dad-9134-856defad3f8e', 8, 'admin', 'Procurement', 'SHIP', 'PurchaseOrder', '44', '{\"after\": {\"id\": 44, \"remark\": \"\", \"status\": \"SHIPPED\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601221716443200\", \"created_at\": \"2026-01-22T17:16:45+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-22T17:17:12+08:00\", \"shipped_at\": \"2026-01-22T17:17:14.3538153+08:00\", \"updated_at\": \"2026-01-22T17:17:14.354+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-22 17:17:14', '2026-01-22 17:17:14');
INSERT INTO `audit_log` VALUES (49, '69f9475f-33ad-4d6f-ac16-34db25c4ddc2', 8, 'admin', 'Procurement', 'RECEIVE', 'PurchaseOrder', '44', '{\"after\": {\"id\": 44, \"items\": [{\"id\": 33, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-22T17:16:45+08:00\", \"updated_at\": \"2026-01-22T17:16:45+08:00\", \"qty_ordered\": 100, \"qty_received\": 100, \"purchase_order_id\": 44}, {\"id\": 34, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-22T17:16:45+08:00\", \"updated_at\": \"2026-01-22T17:16:45+08:00\", \"qty_ordered\": 1, \"qty_received\": 1, \"purchase_order_id\": 44}], \"remark\": \"\", \"status\": \"RECEIVED\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601221716443200\", \"created_at\": \"2026-01-22T17:16:45+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-22T17:17:12+08:00\", \"shipped_at\": \"2026-01-22T17:17:14+08:00\", \"updated_at\": \"2026-01-22T17:17:29.138+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": \"2026-01-22T17:17:29.1271264+08:00\", \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-22 17:17:29', '2026-01-22 17:17:29');
INSERT INTO `audit_log` VALUES (50, 'e8947768-6777-4e92-92d8-abaab105ce5f', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '46', '{\"after\": {\"id\": 46, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601231357241900\", \"created_at\": \"2026-01-23T13:57:24.987+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T13:57:24.987+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 13:57:25', '2026-01-23 13:57:25');
INSERT INTO `audit_log` VALUES (51, '121d66cf-9653-4e48-a178-2def0cc6f56a', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '48', '{\"after\": {\"id\": 48, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601231357340400\", \"created_at\": \"2026-01-23T13:57:34.758+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T13:57:34.758+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 13:57:35', '2026-01-23 13:57:35');
INSERT INTO `audit_log` VALUES (52, '7f224635-4a48-4b1a-8efc-5a44e602b902', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '50', '{\"after\": {\"id\": 50, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601231402361800\", \"created_at\": \"2026-01-23T14:02:36.546+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:02:36.546+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:02:37', '2026-01-23 14:02:37');
INSERT INTO `audit_log` VALUES (53, 'f2dcdd2f-01e9-4abb-8697-488951169857', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '52', '{\"after\": {\"id\": 52, \"items\": [{\"id\": 0, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601231403353900\", \"created_at\": \"2026-01-23T14:03:35.342+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:03:35.342+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:03:35', '2026-01-23 14:03:35');
INSERT INTO `audit_log` VALUES (54, '1e69baeb-1e34-4a66-96da-db03309b1a47', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '53', '{\"after\": {\"id\": 53, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601231403354100\", \"created_at\": \"2026-01-23T14:03:35.342+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:03:35.342+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:03:35', '2026-01-23 14:03:35');
INSERT INTO `audit_log` VALUES (55, '1645f335-30a1-4b40-a939-a81f5449b495', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '54', '{\"after\": {\"id\": 54, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601231403479500\", \"created_at\": \"2026-01-23T14:03:47.673+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:03:47.673+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:03:48', '2026-01-23 14:03:48');
INSERT INTO `audit_log` VALUES (56, '593811d0-3159-4600-89d2-951f5f587c5b', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '56', '{\"after\": {\"id\": 56, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601231404521500\", \"created_at\": \"2026-01-23T14:04:52.109+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:04:52.109+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:04:52', '2026-01-23 14:04:52');
INSERT INTO `audit_log` VALUES (57, '74f2a7f4-443e-4cba-aa4f-5ab83e413ee6', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '58', '{\"after\": {\"id\": 58, \"items\": [{\"id\": 0, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601231406493400\", \"created_at\": \"2026-01-23T14:06:49.612+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:06:49.612+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:06:50', '2026-01-23 14:06:50');
INSERT INTO `audit_log` VALUES (58, '570ec691-e1c5-424d-9836-55771aae44e6', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '59', '{\"after\": {\"id\": 59, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601231406493200\", \"created_at\": \"2026-01-23T14:06:49.613+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:06:49.613+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:06:50', '2026-01-23 14:06:50');
INSERT INTO `audit_log` VALUES (59, '4ddef9f0-4277-446d-9dbb-68297874806a', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '60', '{\"after\": {\"id\": 60, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601231406572800\", \"created_at\": \"2026-01-23T14:06:57.971+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:06:57.971+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:06:58', '2026-01-23 14:06:58');
INSERT INTO `audit_log` VALUES (60, 'afee8683-c204-4039-904b-cc91480cf994', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '62', '{\"after\": {\"id\": 62, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601231408026700\", \"created_at\": \"2026-01-23T14:08:02.739+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:08:02.739+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:08:03', '2026-01-23 14:08:03');
INSERT INTO `audit_log` VALUES (61, 'e93867d6-4799-4b89-b694-a6297f711a94', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '64', '{\"after\": {\"id\": 64, \"items\": [{\"id\": 0, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO20260123140830491000\", \"created_at\": \"2026-01-23T14:08:30.423+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:08:30.423+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:08:30', '2026-01-23 14:08:30');
INSERT INTO `audit_log` VALUES (62, 'ecca949e-d9f0-4417-bf71-53d9157a9fb4', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '66', '{\"after\": {\"id\": 66, \"items\": [{\"id\": 0, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO20260123140845749600\", \"created_at\": \"2026-01-23T14:08:45.442+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:08:45.442+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:08:45', '2026-01-23 14:08:45');
INSERT INTO `audit_log` VALUES (63, '90df3002-b3ed-4a7c-9969-933492af6a1e', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '69', '{\"after\": {\"id\": 69, \"items\": [{\"id\": 0, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO20260123140948677600\", \"created_at\": \"2026-01-23T14:09:48.729+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:09:48.729+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:09:49', '2026-01-23 14:09:49');
INSERT INTO `audit_log` VALUES (64, 'ce1bc7c9-ef98-4432-9016-7926e5860d07', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '68', '{\"after\": {\"id\": 68, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO20260123140948199200\", \"created_at\": \"2026-01-23T14:09:48.729+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:09:48.729+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:09:49', '2026-01-23 14:09:49');
INSERT INTO `audit_log` VALUES (65, '34b421e3-2cb8-449c-aae2-5526d0c44f3d', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '70', '{\"after\": {\"id\": 70, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO20260123140957347500\", \"created_at\": \"2026-01-23T14:09:57.041+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:09:57.041+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:09:57', '2026-01-23 14:09:57');
INSERT INTO `audit_log` VALUES (66, '764d77d3-65e2-4e16-afb1-26e1dfb25ea7', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '72', '{\"after\": {\"id\": 72, \"items\": [{\"id\": 0, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO20260123141031583600\", \"created_at\": \"2026-01-23T14:10:31.371+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:10:31.371+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:10:31', '2026-01-23 14:10:31');
INSERT INTO `audit_log` VALUES (67, '1c7a9dea-9dc2-4b9f-9cab-72277f5ab721', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '74', '{\"after\": {\"id\": 74, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO20260123141349795900\", \"created_at\": \"2026-01-23T14:13:49.039+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:13:49.039+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:13:49', '2026-01-23 14:13:49');
INSERT INTO `audit_log` VALUES (68, '1596a9e0-ac64-4d6c-a130-217cdfecb21b', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '75', '{\"after\": {\"id\": 75, \"items\": [{\"id\": 0, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO20260123141349409900\", \"created_at\": \"2026-01-23T14:13:49.039+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:13:49.039+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:13:49', '2026-01-23 14:13:49');
INSERT INTO `audit_log` VALUES (69, '13f107d6-8d41-4658-af5c-6c94c99e38a3', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '76', '{\"after\": {\"id\": 76, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO20260123141400316900\", \"created_at\": \"2026-01-23T14:14:00.323+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:14:00.323+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:14:00', '2026-01-23 14:14:00');
INSERT INTO `audit_log` VALUES (70, 'e1cafcc8-5dc9-4754-b1f8-3fa92ee77eb4', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '78', '{\"after\": {\"id\": 78, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO20260123141609864d37d4\", \"created_at\": \"2026-01-23T14:16:09.081+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:16:09.081+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:16:09', '2026-01-23 14:16:09');
INSERT INTO `audit_log` VALUES (71, 'cedfa6e2-b039-49d1-9bab-83a1ee716b5b', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '79', '{\"after\": {\"id\": 79, \"items\": [{\"id\": 0, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601231416094da24355\", \"created_at\": \"2026-01-23T14:16:09.088+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:16:09.088+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:16:09', '2026-01-23 14:16:09');
INSERT INTO `audit_log` VALUES (72, 'bdb2f8c2-6a14-4581-90d1-cfe9989ce007', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '80', '{\"after\": {\"id\": 80, \"items\": [{\"id\": 0, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO20260123141628c78db3ce\", \"created_at\": \"2026-01-23T14:16:28.667+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:16:28.667+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:16:29', '2026-01-23 14:16:29');
INSERT INTO `audit_log` VALUES (73, '646725b7-57e9-47f4-8422-c617bcb5d48e', 8, 'admin', 'Procurement', 'CREATE', 'PurchaseOrder', '81', '{\"after\": {\"id\": 81, \"items\": [{\"id\": 0, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 0}, {\"id\": 0, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"0001-01-01T00:00:00Z\", \"updated_at\": \"0001-01-01T00:00:00Z\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 0}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"po_number\": \"PO202601231416282515eabc\", \"created_at\": \"2026-01-23T14:16:28.667+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:16:28.667+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:16:29', '2026-01-23 14:16:29');
INSERT INTO `audit_log` VALUES (74, 'e0d5059e-9080-4a5c-b418-eb4e533244c8', 8, 'admin', 'Procurement', 'SUBMIT', 'PurchaseOrder', '80', '{\"after\": {\"id\": 80, \"remark\": \"\", \"status\": \"ORDERED\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO20260123141628c78db3ce\", \"created_at\": \"2026-01-23T14:16:29+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-23T14:17:33.8014874+08:00\", \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:17:33.803+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:17:34', '2026-01-23 14:17:34');
INSERT INTO `audit_log` VALUES (75, 'ecd3c3bb-1750-47e9-bb1f-3671d317da8e', 8, 'admin', 'Procurement', 'SHIP', 'PurchaseOrder', '80', '{\"after\": {\"id\": 80, \"remark\": \"\", \"status\": \"SHIPPED\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO20260123141628c78db3ce\", \"created_at\": \"2026-01-23T14:16:29+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-23T14:17:34+08:00\", \"shipped_at\": \"2026-01-23T14:17:36.0705327+08:00\", \"updated_at\": \"2026-01-23T14:17:36.071+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:17:36', '2026-01-23 14:17:36');
INSERT INTO `audit_log` VALUES (76, 'dd9af9c6-ba7f-49ea-b6f7-6c03cf5e5ba0', 8, 'admin', 'Procurement', 'RECEIVE', 'PurchaseOrder', '80', '{\"after\": {\"id\": 80, \"items\": [{\"id\": 71, \"sku\": {\"id\": 3, \"title\": \"组合子产品B\", \"seller_sku\": \"SKU-CHILD-003\"}, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:16:29+08:00\", \"updated_at\": \"2026-01-23T14:16:29+08:00\", \"qty_ordered\": 1, \"qty_received\": 1, \"purchase_order_id\": 80}], \"remark\": \"\", \"status\": \"RECEIVED\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO20260123141628c78db3ce\", \"created_at\": \"2026-01-23T14:16:29+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-23T14:17:34+08:00\", \"shipped_at\": \"2026-01-23T14:17:36+08:00\", \"updated_at\": \"2026-01-23T14:17:43.597+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": \"2026-01-23T14:17:43.5879908+08:00\", \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:17:44', '2026-01-23 14:17:44');
INSERT INTO `audit_log` VALUES (77, 'c56fd50e-f5b7-4995-b116-0a6f82826f1b', 8, 'admin', 'Procurement', 'SUBMIT', 'PurchaseOrder', '81', '{\"after\": {\"id\": 81, \"remark\": \"\", \"status\": \"ORDERED\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601231416282515eabc\", \"created_at\": \"2026-01-23T14:16:29+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-23T14:17:54.4703466+08:00\", \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:17:54.47+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:17:54', '2026-01-23 14:17:54');
INSERT INTO `audit_log` VALUES (78, '5a89bea7-5365-45b3-a21c-1f1d79c7fa36', 8, 'admin', 'Procurement', 'SHIP', 'PurchaseOrder', '81', '{\"after\": {\"id\": 81, \"remark\": \"\", \"status\": \"SHIPPED\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601231416282515eabc\", \"created_at\": \"2026-01-23T14:16:29+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-23T14:17:54+08:00\", \"shipped_at\": \"2026-01-23T14:17:57.4228335+08:00\", \"updated_at\": \"2026-01-23T14:17:57.423+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:17:57', '2026-01-23 14:17:57');
INSERT INTO `audit_log` VALUES (79, 'ff726c98-02fa-458c-8722-3c235a35df13', 8, 'admin', 'Procurement', 'RECEIVE', 'PurchaseOrder', '81', '{\"after\": {\"id\": 81, \"items\": [{\"id\": 72, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-23T14:16:29+08:00\", \"updated_at\": \"2026-01-23T14:16:29+08:00\", \"qty_ordered\": 100, \"qty_received\": 100, \"purchase_order_id\": 81}, {\"id\": 73, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:16:29+08:00\", \"updated_at\": \"2026-01-23T14:16:29+08:00\", \"qty_ordered\": 1, \"qty_received\": 1, \"purchase_order_id\": 81}], \"remark\": \"\", \"status\": \"RECEIVED\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601231416282515eabc\", \"created_at\": \"2026-01-23T14:16:29+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-23T14:17:54+08:00\", \"shipped_at\": \"2026-01-23T14:17:57+08:00\", \"updated_at\": \"2026-01-23T14:18:08.798+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": \"2026-01-23T14:18:08.7800662+08:00\", \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-23 14:18:09', '2026-01-23 14:18:09');
INSERT INTO `audit_log` VALUES (80, 'ea54d682-508f-4e22-9a3c-0ffc7e23e39a', 8, 'admin', 'Procurement', 'SUBMIT', 'PurchaseOrder', '79', '{\"after\": {\"id\": 79, \"remark\": \"\", \"status\": \"ORDERED\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO202601231416094da24355\", \"created_at\": \"2026-01-23T14:16:09+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-24T12:26:50.0326663+08:00\", \"shipped_at\": null, \"updated_at\": \"2026-01-24T12:26:50.033+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-24 12:26:50', '2026-01-24 12:26:50');
INSERT INTO `audit_log` VALUES (81, 'e907f6be-ea3e-4bcf-8bb2-eda4bc2394e4', 8, 'admin', 'Procurement', 'SHIP', 'PurchaseOrder', '79', '{\"after\": {\"id\": 79, \"remark\": \"\", \"status\": \"SHIPPED\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO202601231416094da24355\", \"created_at\": \"2026-01-23T14:16:09+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-24T12:26:50+08:00\", \"shipped_at\": \"2026-01-24T12:26:55.5895797+08:00\", \"updated_at\": \"2026-01-24T12:26:55.601+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-24 12:26:56', '2026-01-24 12:26:56');
INSERT INTO `audit_log` VALUES (82, '9b046414-c47c-4c75-8553-c85adc8ce84b', 8, 'admin', 'Procurement', 'RECEIVE', 'PurchaseOrder', '79', '{\"after\": {\"id\": 79, \"items\": [{\"id\": 70, \"sku\": {\"id\": 3, \"title\": \"组合子产品B\", \"seller_sku\": \"SKU-CHILD-003\"}, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:16:09+08:00\", \"updated_at\": \"2026-01-23T14:16:09+08:00\", \"qty_ordered\": 1, \"qty_received\": 1, \"purchase_order_id\": 79}], \"remark\": \"\", \"status\": \"RECEIVED\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO202601231416094da24355\", \"created_at\": \"2026-01-23T14:16:09+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-24T12:26:50+08:00\", \"shipped_at\": \"2026-01-24T12:26:56+08:00\", \"updated_at\": \"2026-01-24T12:27:01.025+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": \"2026-01-24T12:27:01.018433+08:00\", \"supplier_id\": 2, \"total_amount\": 1}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-24 12:27:01', '2026-01-24 12:27:01');
INSERT INTO `audit_log` VALUES (83, '8e760b58-8706-4b65-af11-538080170266', 8, 'admin', 'Procurement', 'SUBMIT', 'PurchaseOrder', '78', '{\"after\": {\"id\": 78, \"remark\": \"\", \"status\": \"ORDERED\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO20260123141609864d37d4\", \"created_at\": \"2026-01-23T14:16:09+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-24T12:27:03.167201+08:00\", \"shipped_at\": null, \"updated_at\": \"2026-01-24T12:27:03.167+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-24 12:27:03', '2026-01-24 12:27:03');
INSERT INTO `audit_log` VALUES (84, 'c29a6c2b-24b8-453a-96b8-36e07928e05a', 8, 'admin', 'Procurement', 'SHIP', 'PurchaseOrder', '78', '{\"after\": {\"id\": 78, \"remark\": \"\", \"status\": \"SHIPPED\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO20260123141609864d37d4\", \"created_at\": \"2026-01-23T14:16:09+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-24T12:27:03+08:00\", \"shipped_at\": \"2026-01-24T12:27:07.7366705+08:00\", \"updated_at\": \"2026-01-24T12:27:07.749+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-24 12:27:08', '2026-01-24 12:27:08');
INSERT INTO `audit_log` VALUES (85, '7132b382-6b2a-4868-9c13-9c66bbadb44e', 8, 'admin', 'Procurement', 'RECEIVE', 'PurchaseOrder', '78', '{\"after\": {\"id\": 78, \"items\": [{\"id\": 68, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-23T14:16:09+08:00\", \"updated_at\": \"2026-01-23T14:16:09+08:00\", \"qty_ordered\": 100, \"qty_received\": 100, \"purchase_order_id\": 78}, {\"id\": 69, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:16:09+08:00\", \"updated_at\": \"2026-01-23T14:16:09+08:00\", \"qty_ordered\": 1, \"qty_received\": 1, \"purchase_order_id\": 78}], \"remark\": \"\", \"status\": \"RECEIVED\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO20260123141609864d37d4\", \"created_at\": \"2026-01-23T14:16:09+08:00\", \"created_by\": 8, \"ordered_at\": \"2026-01-24T12:27:03+08:00\", \"shipped_at\": \"2026-01-24T12:27:08+08:00\", \"updated_at\": \"2026-01-24T12:27:12.572+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": \"2026-01-24T12:27:12.5604083+08:00\", \"supplier_id\": 1, \"total_amount\": 10001}, \"before\": null}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-24 12:27:13', '2026-01-24 12:27:13');
INSERT INTO `audit_log` VALUES (86, '745f45e8-e99a-418d-9266-d37d01da8597', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '52', '{\"after\": null, \"before\": {\"id\": 52, \"items\": [{\"id\": 41, \"sku\": {\"id\": 3, \"title\": \"组合子产品B\", \"seller_sku\": \"SKU-CHILD-003\"}, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:03:35+08:00\", \"updated_at\": \"2026-01-23T14:03:35+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 52}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO202601231403353900\", \"created_at\": \"2026-01-23T14:03:35+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:03:35+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:54:36', '2026-01-25 20:54:36');
INSERT INTO `audit_log` VALUES (87, 'd33c168f-0378-4f3d-9f61-1fc9a0f1df05', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '53', '{\"after\": null, \"before\": {\"id\": 53, \"items\": [{\"id\": 42, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-23T14:03:35+08:00\", \"updated_at\": \"2026-01-23T14:03:35+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 53}, {\"id\": 43, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:03:35+08:00\", \"updated_at\": \"2026-01-23T14:03:35+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 53}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601231403354100\", \"created_at\": \"2026-01-23T14:03:35+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:03:35+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:54:39', '2026-01-25 20:54:39');
INSERT INTO `audit_log` VALUES (88, '8d871ff1-2bb7-48ba-bada-47b6a16fb781', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '50', '{\"after\": null, \"before\": {\"id\": 50, \"items\": [{\"id\": 39, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-23T14:02:37+08:00\", \"updated_at\": \"2026-01-23T14:02:37+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 50}, {\"id\": 40, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:02:37+08:00\", \"updated_at\": \"2026-01-23T14:02:37+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 50}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601231402361800\", \"created_at\": \"2026-01-23T14:02:37+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:02:37+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:54:42', '2026-01-25 20:54:42');
INSERT INTO `audit_log` VALUES (89, '387d819a-913f-483a-9a67-d1b622b78103', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '48', '{\"after\": null, \"before\": {\"id\": 48, \"items\": [{\"id\": 37, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-23T13:57:35+08:00\", \"updated_at\": \"2026-01-23T13:57:35+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 48}, {\"id\": 38, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T13:57:35+08:00\", \"updated_at\": \"2026-01-23T13:57:35+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 48}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601231357340400\", \"created_at\": \"2026-01-23T13:57:35+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T13:57:35+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:54:43', '2026-01-25 20:54:43');
INSERT INTO `audit_log` VALUES (90, '3049cbb2-9682-484c-b4bc-462a0ec614a6', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '46', '{\"after\": null, \"before\": {\"id\": 46, \"items\": [{\"id\": 35, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-23T13:57:25+08:00\", \"updated_at\": \"2026-01-23T13:57:25+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 46}, {\"id\": 36, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T13:57:25+08:00\", \"updated_at\": \"2026-01-23T13:57:25+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 46}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601231357241900\", \"created_at\": \"2026-01-23T13:57:25+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T13:57:25+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:54:46', '2026-01-25 20:54:46');
INSERT INTO `audit_log` VALUES (91, '8d6b360d-ee45-40f4-bf3f-e77d95d9cde7', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '54', '{\"after\": null, \"before\": {\"id\": 54, \"items\": [{\"id\": 44, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-23T14:03:48+08:00\", \"updated_at\": \"2026-01-23T14:03:48+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 54}, {\"id\": 45, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:03:48+08:00\", \"updated_at\": \"2026-01-23T14:03:48+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 54}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601231403479500\", \"created_at\": \"2026-01-23T14:03:48+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:03:48+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:54:48', '2026-01-25 20:54:48');
INSERT INTO `audit_log` VALUES (92, '5739db56-9a6b-495a-a951-a4b4579bebbe', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '42', '{\"after\": null, \"before\": {\"id\": 42, \"items\": [{\"id\": 31, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-22T17:16:36+08:00\", \"updated_at\": \"2026-01-22T17:16:36+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 42}, {\"id\": 32, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-22T17:16:36+08:00\", \"updated_at\": \"2026-01-22T17:16:36+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 42}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601221716361900\", \"created_at\": \"2026-01-22T17:16:36+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-22T17:16:36+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:54:51', '2026-01-25 20:54:51');
INSERT INTO `audit_log` VALUES (93, '326ac625-05d3-4f6a-a158-5144ee217dd1', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '76', '{\"after\": null, \"before\": {\"id\": 76, \"items\": [{\"id\": 66, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-23T14:14:00+08:00\", \"updated_at\": \"2026-01-23T14:14:00+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 76}, {\"id\": 67, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:14:00+08:00\", \"updated_at\": \"2026-01-23T14:14:00+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 76}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO20260123141400316900\", \"created_at\": \"2026-01-23T14:14:00+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:14:00+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:54:55', '2026-01-25 20:54:55');
INSERT INTO `audit_log` VALUES (94, 'ccc22744-23d5-4e73-a1ed-26fc2a10480a', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '75', '{\"after\": null, \"before\": {\"id\": 75, \"items\": [{\"id\": 65, \"sku\": {\"id\": 3, \"title\": \"组合子产品B\", \"seller_sku\": \"SKU-CHILD-003\"}, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:13:49+08:00\", \"updated_at\": \"2026-01-23T14:13:49+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 75}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO20260123141349409900\", \"created_at\": \"2026-01-23T14:13:49+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:13:49+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:54:57', '2026-01-25 20:54:57');
INSERT INTO `audit_log` VALUES (95, 'bd11fcd4-777f-4234-b797-0ce1e4bae449', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '74', '{\"after\": null, \"before\": {\"id\": 74, \"items\": [{\"id\": 63, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-23T14:13:49+08:00\", \"updated_at\": \"2026-01-23T14:13:49+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 74}, {\"id\": 64, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:13:49+08:00\", \"updated_at\": \"2026-01-23T14:13:49+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 74}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO20260123141349795900\", \"created_at\": \"2026-01-23T14:13:49+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:13:49+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:54:58', '2026-01-25 20:54:58');
INSERT INTO `audit_log` VALUES (96, '5fdd8568-36ce-4fe3-8a4f-26669b9c0321', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '72', '{\"after\": null, \"before\": {\"id\": 72, \"items\": [{\"id\": 62, \"sku\": {\"id\": 3, \"title\": \"组合子产品B\", \"seller_sku\": \"SKU-CHILD-003\"}, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:10:31+08:00\", \"updated_at\": \"2026-01-23T14:10:31+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 72}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO20260123141031583600\", \"created_at\": \"2026-01-23T14:10:31+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:10:31+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:01', '2026-01-25 20:55:01');
INSERT INTO `audit_log` VALUES (97, 'f8282a48-6a98-4c72-bc5c-22ee16c07ccc', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '70', '{\"after\": null, \"before\": {\"id\": 70, \"items\": [{\"id\": 60, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-23T14:09:57+08:00\", \"updated_at\": \"2026-01-23T14:09:57+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 70}, {\"id\": 61, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:09:57+08:00\", \"updated_at\": \"2026-01-23T14:09:57+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 70}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO20260123140957347500\", \"created_at\": \"2026-01-23T14:09:57+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:09:57+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:03', '2026-01-25 20:55:03');
INSERT INTO `audit_log` VALUES (98, 'a8e0de60-016c-457d-bb1d-7024894843a6', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '68', '{\"after\": null, \"before\": {\"id\": 68, \"items\": [{\"id\": 58, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-23T14:09:49+08:00\", \"updated_at\": \"2026-01-23T14:09:49+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 68}, {\"id\": 59, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:09:49+08:00\", \"updated_at\": \"2026-01-23T14:09:49+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 68}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO20260123140948199200\", \"created_at\": \"2026-01-23T14:09:49+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:09:49+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:05', '2026-01-25 20:55:05');
INSERT INTO `audit_log` VALUES (99, '28a43cc7-00f6-4992-9125-79bb5ead14d0', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '69', '{\"after\": null, \"before\": {\"id\": 69, \"items\": [{\"id\": 57, \"sku\": {\"id\": 3, \"title\": \"组合子产品B\", \"seller_sku\": \"SKU-CHILD-003\"}, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:09:49+08:00\", \"updated_at\": \"2026-01-23T14:09:49+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 69}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO20260123140948677600\", \"created_at\": \"2026-01-23T14:09:49+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:09:49+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:06', '2026-01-25 20:55:06');
INSERT INTO `audit_log` VALUES (100, 'd207f9cb-c971-4382-b0af-9c1113cf1e27', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '66', '{\"after\": null, \"before\": {\"id\": 66, \"items\": [{\"id\": 56, \"sku\": {\"id\": 3, \"title\": \"组合子产品B\", \"seller_sku\": \"SKU-CHILD-003\"}, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:08:45+08:00\", \"updated_at\": \"2026-01-23T14:08:45+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 66}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO20260123140845749600\", \"created_at\": \"2026-01-23T14:08:45+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:08:45+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:07', '2026-01-25 20:55:07');
INSERT INTO `audit_log` VALUES (101, '76768024-9c03-474d-8701-207d1c54cba7', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '64', '{\"after\": null, \"before\": {\"id\": 64, \"items\": [{\"id\": 55, \"sku\": {\"id\": 3, \"title\": \"组合子产品B\", \"seller_sku\": \"SKU-CHILD-003\"}, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:08:30+08:00\", \"updated_at\": \"2026-01-23T14:08:30+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 64}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO20260123140830491000\", \"created_at\": \"2026-01-23T14:08:30+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:08:30+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:09', '2026-01-25 20:55:09');
INSERT INTO `audit_log` VALUES (102, 'ccd8385b-a0c2-445e-a2e9-8e2a35c180ad', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '62', '{\"after\": null, \"before\": {\"id\": 62, \"items\": [{\"id\": 53, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-23T14:08:03+08:00\", \"updated_at\": \"2026-01-23T14:08:03+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 62}, {\"id\": 54, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:08:03+08:00\", \"updated_at\": \"2026-01-23T14:08:03+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 62}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601231408026700\", \"created_at\": \"2026-01-23T14:08:03+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:08:03+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:10', '2026-01-25 20:55:10');
INSERT INTO `audit_log` VALUES (103, '928a7a8c-9cfb-466e-b18c-2a1f23b8218c', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '60', '{\"after\": null, \"before\": {\"id\": 60, \"items\": [{\"id\": 51, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-23T14:06:58+08:00\", \"updated_at\": \"2026-01-23T14:06:58+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 60}, {\"id\": 52, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:06:58+08:00\", \"updated_at\": \"2026-01-23T14:06:58+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 60}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601231406572800\", \"created_at\": \"2026-01-23T14:06:58+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:06:58+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:12', '2026-01-25 20:55:12');
INSERT INTO `audit_log` VALUES (104, '81d5b39a-d10e-4d5f-9ed8-322dbdc1f4e1', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '58', '{\"after\": null, \"before\": {\"id\": 58, \"items\": [{\"id\": 48, \"sku\": {\"id\": 3, \"title\": \"组合子产品B\", \"seller_sku\": \"SKU-CHILD-003\"}, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:06:50+08:00\", \"updated_at\": \"2026-01-23T14:06:50+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 58}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO202601231406493400\", \"created_at\": \"2026-01-23T14:06:50+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:06:50+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:14', '2026-01-25 20:55:14');
INSERT INTO `audit_log` VALUES (105, 'bc209f8c-5803-48bb-938a-c7d01657ef10', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '59', '{\"after\": null, \"before\": {\"id\": 59, \"items\": [{\"id\": 49, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-23T14:06:50+08:00\", \"updated_at\": \"2026-01-23T14:06:50+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 59}, {\"id\": 50, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:06:50+08:00\", \"updated_at\": \"2026-01-23T14:06:50+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 59}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601231406493200\", \"created_at\": \"2026-01-23T14:06:50+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:06:50+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:15', '2026-01-25 20:55:15');
INSERT INTO `audit_log` VALUES (106, '3f2cb8e4-bcca-44af-9630-1eebe83b0889', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '56', '{\"after\": null, \"before\": {\"id\": 56, \"items\": [{\"id\": 46, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-23T14:04:52+08:00\", \"updated_at\": \"2026-01-23T14:04:52+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 56}, {\"id\": 47, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-23T14:04:52+08:00\", \"updated_at\": \"2026-01-23T14:04:52+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 56}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601231404521500\", \"created_at\": \"2026-01-23T14:04:52+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-23T14:04:52+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:17', '2026-01-25 20:55:17');
INSERT INTO `audit_log` VALUES (107, '6b17ceaa-0df3-46c1-b811-97a1c4af652e', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '41', '{\"after\": null, \"before\": {\"id\": 41, \"items\": [{\"id\": 29, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-21T20:00:50+08:00\", \"updated_at\": \"2026-01-21T20:00:50+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 41}, {\"id\": 30, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-21T20:00:50+08:00\", \"updated_at\": \"2026-01-21T20:00:50+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 41}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601212000499100\", \"created_at\": \"2026-01-21T20:00:50+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T20:00:50+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:19', '2026-01-25 20:55:19');
INSERT INTO `audit_log` VALUES (108, '79e09452-e511-434f-88df-9ff328db3d7a', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '38', '{\"after\": null, \"before\": {\"id\": 38, \"items\": [{\"id\": 27, \"sku\": {\"id\": 3, \"title\": \"组合子产品B\", \"seller_sku\": \"SKU-CHILD-003\"}, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-21T19:59:31+08:00\", \"updated_at\": \"2026-01-21T19:59:31+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 38}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO202601211959315700\", \"created_at\": \"2026-01-21T19:59:31+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T19:59:31+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:21', '2026-01-25 20:55:21');
INSERT INTO `audit_log` VALUES (109, 'a395cc44-f473-4e0e-aa53-1bb482e7e3ff', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '37', '{\"after\": null, \"before\": {\"id\": 37, \"items\": [{\"id\": 24, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 199.9, \"unit_cost\": 19.99, \"created_at\": \"2026-01-21T19:57:08+08:00\", \"updated_at\": \"2026-01-21T19:57:08+08:00\", \"qty_ordered\": 10, \"qty_received\": 0, \"purchase_order_id\": 37}, {\"id\": 25, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 99.9, \"unit_cost\": 9.99, \"created_at\": \"2026-01-21T19:57:08+08:00\", \"updated_at\": \"2026-01-21T19:57:08+08:00\", \"qty_ordered\": 10, \"qty_received\": 0, \"purchase_order_id\": 37}, {\"id\": 26, \"sku\": {\"id\": 3, \"title\": \"组合子产品B\", \"seller_sku\": \"SKU-CHILD-003\"}, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 129.9, \"unit_cost\": 12.99, \"created_at\": \"2026-01-21T19:57:08+08:00\", \"updated_at\": \"2026-01-21T19:57:08+08:00\", \"qty_ordered\": 10, \"qty_received\": 0, \"purchase_order_id\": 37}], \"remark\": \"���Բɹ���\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 4, \"name\": \"广州电子元件\"}, \"po_number\": \"PO202601211957070800\", \"created_at\": \"2026-01-21T19:57:08+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T19:57:08+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 4, \"total_amount\": 429.7}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:23', '2026-01-25 20:55:23');
INSERT INTO `audit_log` VALUES (110, 'cbedf857-e118-4242-9ebc-81be24310bf9', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '35', '{\"after\": null, \"before\": {\"id\": 35, \"items\": [{\"id\": 22, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-21T19:55:14+08:00\", \"updated_at\": \"2026-01-21T19:55:14+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 35}, {\"id\": 23, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-21T19:55:14+08:00\", \"updated_at\": \"2026-01-21T19:55:14+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 35}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601211955149100\", \"created_at\": \"2026-01-21T19:55:14+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T19:55:14+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:24', '2026-01-25 20:55:24');
INSERT INTO `audit_log` VALUES (111, '816594b9-45b2-4ad0-ad78-fc16acf64c2e', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '33', '{\"after\": null, \"before\": {\"id\": 33, \"items\": [{\"id\": 21, \"sku\": {\"id\": 3, \"title\": \"组合子产品B\", \"seller_sku\": \"SKU-CHILD-003\"}, \"sku_id\": 3, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-21T19:55:09+08:00\", \"updated_at\": \"2026-01-21T19:55:09+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 33}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 2, \"name\": \"杭州包装工厂\"}, \"po_number\": \"PO202601211955082500\", \"created_at\": \"2026-01-21T19:55:09+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T19:55:09+08:00\", \"updated_by\": 8, \"marketplace\": \"CA\", \"received_at\": null, \"supplier_id\": 2, \"total_amount\": 1}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:26', '2026-01-25 20:55:26');
INSERT INTO `audit_log` VALUES (112, '569419ff-eeb7-431c-8955-2fc721da3d86', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '25', '{\"after\": null, \"before\": {\"id\": 25, \"items\": [{\"id\": 17, \"sku\": {\"id\": 5, \"title\": \"Test Product 102\", \"seller_sku\": \"SKU-TEST-102\"}, \"sku_id\": 5, \"currency\": \"USD\", \"subtotal\": 3, \"unit_cost\": 3, \"created_at\": \"2026-01-21T19:05:09+08:00\", \"updated_at\": \"2026-01-21T19:05:09+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 25}, {\"id\": 18, \"sku\": {\"id\": 6, \"title\": \"Test Product 103\", \"seller_sku\": \"SKU-TEST-103\"}, \"sku_id\": 6, \"currency\": \"USD\", \"subtotal\": 3, \"unit_cost\": 3, \"created_at\": \"2026-01-21T19:05:09+08:00\", \"updated_at\": \"2026-01-21T19:05:09+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 25}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 3, \"name\": \"上海物流服务\"}, \"po_number\": \"PO202601211905087700\", \"created_at\": \"2026-01-21T19:05:09+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T19:05:09+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 3, \"total_amount\": 6}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:28', '2026-01-25 20:55:28');
INSERT INTO `audit_log` VALUES (113, '77cc0127-0e16-49cc-b912-db4f6701b0a1', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '21', '{\"after\": null, \"before\": {\"id\": 21, \"items\": [{\"id\": 16, \"sku\": {\"id\": 4, \"title\": \"Test Product 101\", \"seller_sku\": \"SKU-TEST-101\"}, \"sku_id\": 4, \"currency\": \"USD\", \"subtotal\": 2, \"unit_cost\": 2, \"created_at\": \"2026-01-21T19:04:57+08:00\", \"updated_at\": \"2026-01-21T19:04:57+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 21}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 4, \"name\": \"广州电子元件\"}, \"po_number\": \"PO202601211904565500\", \"created_at\": \"2026-01-21T19:04:57+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T19:04:57+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 4, \"total_amount\": 2}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:30', '2026-01-25 20:55:30');
INSERT INTO `audit_log` VALUES (114, '50855186-94d5-4812-9439-c3934d46d2f8', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '22', '{\"after\": null, \"before\": {\"id\": 22, \"items\": [{\"id\": 14, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-21T19:04:57+08:00\", \"updated_at\": \"2026-01-21T19:04:57+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 22}, {\"id\": 15, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-21T19:04:57+08:00\", \"updated_at\": \"2026-01-21T19:04:57+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 22}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601211904568100\", \"created_at\": \"2026-01-21T19:04:57+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T19:04:57+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:31', '2026-01-25 20:55:31');
INSERT INTO `audit_log` VALUES (115, '11ad99de-d19d-4818-9e40-75392d61c705', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '17', '{\"after\": null, \"before\": {\"id\": 17, \"items\": [{\"id\": 11, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-21T18:53:22+08:00\", \"updated_at\": \"2026-01-21T18:53:22+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 17}, {\"id\": 12, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-21T18:53:22+08:00\", \"updated_at\": \"2026-01-21T18:53:22+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 17}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601211853227800\", \"created_at\": \"2026-01-21T18:53:22+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T18:53:22+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:32', '2026-01-25 20:55:32');
INSERT INTO `audit_log` VALUES (116, 'f6b881dd-6919-4c6b-9880-6c321d1fc09f', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '19', '{\"after\": null, \"before\": {\"id\": 19, \"items\": [{\"id\": 13, \"sku\": {\"id\": 4, \"title\": \"Test Product 101\", \"seller_sku\": \"SKU-TEST-101\"}, \"sku_id\": 4, \"currency\": \"USD\", \"subtotal\": 2, \"unit_cost\": 2, \"created_at\": \"2026-01-21T18:53:22+08:00\", \"updated_at\": \"2026-01-21T18:53:22+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 19}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 4, \"name\": \"广州电子元件\"}, \"po_number\": \"PO202601211853225000\", \"created_at\": \"2026-01-21T18:53:22+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T18:53:22+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 4, \"total_amount\": 2}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:34', '2026-01-25 20:55:34');
INSERT INTO `audit_log` VALUES (117, '9d8739b6-cc67-4208-ae6c-0726245cd917', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '13', '{\"after\": null, \"before\": {\"id\": 13, \"items\": [{\"id\": 9, \"sku\": {\"id\": 5, \"title\": \"Test Product 102\", \"seller_sku\": \"SKU-TEST-102\"}, \"sku_id\": 5, \"currency\": \"USD\", \"subtotal\": 3, \"unit_cost\": 3, \"created_at\": \"2026-01-21T18:47:54+08:00\", \"updated_at\": \"2026-01-21T18:47:54+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 13}, {\"id\": 10, \"sku\": {\"id\": 6, \"title\": \"Test Product 103\", \"seller_sku\": \"SKU-TEST-103\"}, \"sku_id\": 6, \"currency\": \"USD\", \"subtotal\": 3, \"unit_cost\": 3, \"created_at\": \"2026-01-21T18:47:54+08:00\", \"updated_at\": \"2026-01-21T18:47:54+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 13}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 3, \"name\": \"上海物流服务\"}, \"po_number\": \"PO202601211847530400\", \"created_at\": \"2026-01-21T18:47:54+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T18:47:54+08:00\", \"updated_by\": 8, \"marketplace\": \"US\", \"received_at\": null, \"supplier_id\": 3, \"total_amount\": 6}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:36', '2026-01-25 20:55:36');
INSERT INTO `audit_log` VALUES (118, 'dc8b998a-9361-4036-9c23-3d4d415462b5', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '9', '{\"after\": null, \"before\": {\"id\": 9, \"items\": [{\"id\": 4, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-21T18:47:32+08:00\", \"updated_at\": \"2026-01-21T18:47:32+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 9}, {\"id\": 5, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-21T18:47:32+08:00\", \"updated_at\": \"2026-01-21T18:47:32+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 9}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601211847321100\", \"created_at\": \"2026-01-21T18:47:32+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T18:47:32+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:37', '2026-01-25 20:55:37');
INSERT INTO `audit_log` VALUES (119, 'ba75981d-e489-4f90-aa08-c6d3e272896f', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '11', '{\"after\": null, \"before\": {\"id\": 11, \"items\": [{\"id\": 6, \"sku\": {\"id\": 5, \"title\": \"Test Product 102\", \"seller_sku\": \"SKU-TEST-102\"}, \"sku_id\": 5, \"currency\": \"USD\", \"subtotal\": 3, \"unit_cost\": 3, \"created_at\": \"2026-01-21T18:47:32+08:00\", \"updated_at\": \"2026-01-21T18:47:32+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 11}, {\"id\": 7, \"sku\": {\"id\": 6, \"title\": \"Test Product 103\", \"seller_sku\": \"SKU-TEST-103\"}, \"sku_id\": 6, \"currency\": \"USD\", \"subtotal\": 3, \"unit_cost\": 3, \"created_at\": \"2026-01-21T18:47:32+08:00\", \"updated_at\": \"2026-01-21T18:47:32+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 11}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 3, \"name\": \"上海物流服务\"}, \"po_number\": \"PO202601211847322400\", \"created_at\": \"2026-01-21T18:47:32+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T18:47:32+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 3, \"total_amount\": 6}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:39', '2026-01-25 20:55:39');
INSERT INTO `audit_log` VALUES (120, '7446f2e4-0310-4822-9a9b-656d9bd09e8e', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '12', '{\"after\": null, \"before\": {\"id\": 12, \"items\": [{\"id\": 8, \"sku\": {\"id\": 4, \"title\": \"Test Product 101\", \"seller_sku\": \"SKU-TEST-101\"}, \"sku_id\": 4, \"currency\": \"USD\", \"subtotal\": 2, \"unit_cost\": 2, \"created_at\": \"2026-01-21T18:47:32+08:00\", \"updated_at\": \"2026-01-21T18:47:32+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 12}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 4, \"name\": \"广州电子元件\"}, \"po_number\": \"PO202601211847323400\", \"created_at\": \"2026-01-21T18:47:32+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T18:47:32+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 4, \"total_amount\": 2}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:41', '2026-01-25 20:55:41');
INSERT INTO `audit_log` VALUES (121, '97071e05-5fd9-4468-8cf0-663a070c8eff', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '5', '{\"after\": null, \"before\": {\"id\": 5, \"items\": [{\"id\": 3, \"sku\": {\"id\": 4, \"title\": \"Test Product 101\", \"seller_sku\": \"SKU-TEST-101\"}, \"sku_id\": 4, \"currency\": \"USD\", \"subtotal\": 2, \"unit_cost\": 2, \"created_at\": \"2026-01-21T18:38:29+08:00\", \"updated_at\": \"2026-01-21T18:38:29+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 5}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 4, \"name\": \"广州电子元件\"}, \"po_number\": \"PO202601211838286900\", \"created_at\": \"2026-01-21T18:38:29+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T18:38:29+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 4, \"total_amount\": 2}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:42', '2026-01-25 20:55:42');
INSERT INTO `audit_log` VALUES (122, 'c2d3f647-0f7f-40c8-8741-0e9e1358d438', 8, 'admin', 'Procurement', 'DELETE', 'PurchaseOrder', '1', '{\"after\": null, \"before\": {\"id\": 1, \"items\": [{\"id\": 1, \"sku\": {\"id\": 1, \"title\": \"组合主产品\", \"seller_sku\": \"SKU-MAIN-001\"}, \"sku_id\": 1, \"currency\": \"USD\", \"subtotal\": 10000, \"unit_cost\": 100, \"created_at\": \"2026-01-21T18:38:19+08:00\", \"updated_at\": \"2026-01-21T18:38:19+08:00\", \"qty_ordered\": 100, \"qty_received\": 0, \"purchase_order_id\": 1}, {\"id\": 2, \"sku\": {\"id\": 2, \"title\": \"组合子产品A\", \"seller_sku\": \"SKU-CHILD-002\"}, \"sku_id\": 2, \"currency\": \"USD\", \"subtotal\": 1, \"unit_cost\": 1, \"created_at\": \"2026-01-21T18:38:19+08:00\", \"updated_at\": \"2026-01-21T18:38:19+08:00\", \"qty_ordered\": 1, \"qty_received\": 0, \"purchase_order_id\": 1}], \"remark\": \"\", \"status\": \"DRAFT\", \"currency\": \"USD\", \"supplier\": {\"id\": 1, \"name\": \"深圳优品供应链\"}, \"po_number\": \"PO202601211838181700\", \"created_at\": \"2026-01-21T18:38:19+08:00\", \"created_by\": 8, \"ordered_at\": null, \"shipped_at\": null, \"updated_at\": \"2026-01-21T18:38:19+08:00\", \"updated_by\": 8, \"marketplace\": \"\", \"received_at\": null, \"supplier_id\": 1, \"total_amount\": 10001}}', '::1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36', '2026-01-25 20:55:44', '2026-01-25 20:55:44');

-- ----------------------------
-- Table structure for cash_ledger
-- ----------------------------
DROP TABLE IF EXISTS `cash_ledger`;
CREATE TABLE `cash_ledger`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `trace_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '追踪ID',
  `ledger_type` enum('INCOME','EXPENSE') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '流水类型: INCOME-收入, EXPENSE-支出',
  `category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '类别: SALES_REVENUE(销售收入), PURCHASE_COST(采购成本), SHIPPING_FEE(运费), PACKAGING_COST(包装成本), OTHER_INCOME(其他收入), OTHER_EXPENSE(其他支出)',
  `amount` decimal(15, 2) NOT NULL COMMENT '金额',
  `currency` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'CNY' COMMENT '货币',
  `reference_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '关联单据类型: PURCHASE_ORDER, SHIPMENT, MANUAL等',
  `reference_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '关联单据ID',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '说明',
  `occurred_at` datetime NOT NULL COMMENT '发生日期',
  `created_by` bigint UNSIGNED NOT NULL COMMENT '创建人',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_trace_id`(`trace_id` ASC) USING BTREE,
  INDEX `idx_ledger_type`(`ledger_type` ASC) USING BTREE,
  INDEX `idx_category`(`category` ASC) USING BTREE,
  INDEX `idx_occurred_at`(`occurred_at` ASC) USING BTREE,
  INDEX `idx_reference_type_reference_id`(`reference_type` ASC, `reference_id` ASC) USING BTREE,
  INDEX `idx_created_by`(`created_by` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '现金流水表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cash_ledger
-- ----------------------------

-- ----------------------------
-- Table structure for costing_snapshot
-- ----------------------------
DROP TABLE IF EXISTS `costing_snapshot`;
CREATE TABLE `costing_snapshot`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `trace_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '追踪ID',
  `sku_id` bigint UNSIGNED NOT NULL COMMENT 'SKU ID',
  `cost_type` enum('PURCHASE','LANDED','AVERAGE') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '成本类型: PURCHASE-采购成本, LANDED-到岸成本, AVERAGE-平均成本',
  `unit_cost` decimal(15, 4) NOT NULL COMMENT '单位成本',
  `currency` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'CNY' COMMENT '货币',
  `effective_from` datetime NOT NULL COMMENT '生效日期',
  `effective_to` datetime NULL DEFAULT NULL COMMENT '失效日期 (NULL表示当前有效)',
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '备注',
  `created_by` bigint UNSIGNED NOT NULL COMMENT '创建人',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_trace_id`(`trace_id` ASC) USING BTREE,
  INDEX `idx_sku_id`(`sku_id` ASC) USING BTREE,
  INDEX `idx_cost_type`(`cost_type` ASC) USING BTREE,
  INDEX `idx_sku_cost_from`(`sku_id` ASC, `cost_type` ASC, `effective_from` ASC) USING BTREE,
  INDEX `idx_sku_cost_to`(`sku_id` ASC, `cost_type` ASC, `effective_to` ASC) USING BTREE,
  INDEX `idx_created_by`(`created_by` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '成本快照表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of costing_snapshot
-- ----------------------------

-- ----------------------------
-- Table structure for field_label
-- ----------------------------
DROP TABLE IF EXISTS `field_label`;
CREATE TABLE `field_label`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `label_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文案Key（小写点分层）',
  `module` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '模块标识',
  `scene` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '场景标识',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '状态',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '备注',
  `labels` json NOT NULL COMMENT '多语言文案（key=locale）',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_label_key`(`label_key` ASC) USING BTREE,
  INDEX `idx_label_module`(`module` ASC) USING BTREE,
  INDEX `idx_label_scene`(`scene` ASC) USING BTREE,
  INDEX `idx_label_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 533 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '通用 i18n 文案配置表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of field_label
-- ----------------------------
INSERT INTO `field_label` VALUES (1, 'global.search', NULL, NULL, 'active', NULL, '{\"en-US\": \"Search\", \"zh-CN\": \"查询\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (2, 'global.reset', NULL, NULL, 'active', NULL, '{\"en-US\": \"Reset\", \"zh-CN\": \"重置\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (3, 'global.view', NULL, NULL, 'active', NULL, '{\"en-US\": \"View\", \"zh-CN\": \"查看\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (4, 'global.edit', NULL, NULL, 'active', NULL, '{\"en-US\": \"Edit\", \"zh-CN\": \"编辑\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (5, 'global.delete', NULL, NULL, 'active', NULL, '{\"en-US\": \"Delete\", \"zh-CN\": \"删除\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (6, 'global.save', NULL, NULL, 'active', NULL, '{\"en-US\": \"Save\", \"zh-CN\": \"保存\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (7, 'global.cancel', NULL, NULL, 'active', NULL, '{\"en-US\": \"Cancel\", \"zh-CN\": \"取消\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (8, 'global.confirm', NULL, NULL, 'active', NULL, '{\"en-US\": \"Confirm\", \"zh-CN\": \"确认\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (9, 'global.warning', NULL, NULL, 'active', NULL, '{\"en-US\": \"Warning\", \"zh-CN\": \"提示\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (10, 'global.required', NULL, NULL, 'active', NULL, '{\"en-US\": \"Required\", \"zh-CN\": \"必填\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (11, 'global.actions', NULL, NULL, 'active', NULL, '{\"en-US\": \"Actions\", \"zh-CN\": \"操作\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (12, 'global.add', NULL, NULL, 'active', NULL, '{\"en-US\": \"Add\", \"zh-CN\": \"新增\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (13, 'global.refresh', NULL, NULL, 'active', NULL, '{\"en-US\": \"Refresh\", \"zh-CN\": \"刷新\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (14, 'global.export', NULL, NULL, 'active', NULL, '{\"en-US\": \"Export\", \"zh-CN\": \"导出\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (15, 'global.import', NULL, NULL, 'active', NULL, '{\"en-US\": \"Import\", \"zh-CN\": \"导入\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (16, 'product.list.skuManagement', 'product', 'list', 'active', NULL, '{\"en-US\": \"SKU Management\", \"zh-CN\": \"产品列表\"}', '2026-01-19 15:26:10', '2026-01-19 22:34:30');
INSERT INTO `field_label` VALUES (17, 'product.list.createSku', 'product', 'list', 'active', NULL, '{\"en-US\": \"Create SKU\", \"zh-CN\": \"新增SKU\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (18, 'product.list.keyword', 'product', 'list', 'active', NULL, '{\"en-US\": \"Keyword\", \"zh-CN\": \"关键词\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (19, 'product.list.keywordPlaceholder', 'product', 'list', 'active', NULL, '{\"en-US\": \"Seller SKU, ASIN or Title\", \"zh-CN\": \"Seller SKU、ASIN 或标题\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (20, 'product.list.combo', 'product', 'list', 'active', NULL, '{\"en-US\": \"Combo\", \"zh-CN\": \"组合\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (21, 'product.list.all', 'product', 'list', 'active', NULL, '{\"en-US\": \"All\", \"zh-CN\": \"全部\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (22, 'product.list.mainOnly', 'product', 'list', 'active', NULL, '{\"en-US\": \"Main Only\", \"zh-CN\": \"仅主SKU\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (23, 'product.list.comboOnly', 'product', 'list', 'active', NULL, '{\"en-US\": \"Combo Only\", \"zh-CN\": \"组合SKU\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (24, 'product.list.marketplace', 'product', 'list', 'active', NULL, '{\"en-US\": \"Marketplace\", \"zh-CN\": \"站点\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (25, 'product.list.status', 'product', 'list', 'active', NULL, '{\"en-US\": \"Status\", \"zh-CN\": \"状态\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (26, 'product.list.statusActive', 'product', 'list', 'active', NULL, '{\"en-US\": \"Active\", \"zh-CN\": \"启用\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (27, 'product.list.statusInactive', 'product', 'list', 'active', NULL, '{\"en-US\": \"Inactive\", \"zh-CN\": \"停用\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (28, 'product.list.statusDiscontinued', 'product', 'list', 'active', NULL, '{\"en-US\": \"Discontinued\", \"zh-CN\": \"停售\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (29, 'product.list.id', 'product', 'list', 'active', NULL, '{\"en-US\": \"ID\", \"zh-CN\": \"ID\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (30, 'product.list.image', 'product', 'list', 'active', NULL, '{\"en-US\": \"Image\", \"zh-CN\": \"图片\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (31, 'product.list.noImage', 'product', 'list', 'active', NULL, '{\"en-US\": \"No Image\", \"zh-CN\": \"暂无图片\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (32, 'product.list.sellerSku', 'product', 'list', 'active', NULL, '{\"en-US\": \"Seller SKU\", \"zh-CN\": \"Seller SKU\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (33, 'product.list.asin', 'product', 'list', 'active', NULL, '{\"en-US\": \"ASIN\", \"zh-CN\": \"ASIN\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (34, 'product.list.title', 'product', 'list', 'active', NULL, '{\"en-US\": \"Title\", \"zh-CN\": \"标题\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (35, 'product.list.supplier', 'product', 'list', 'active', NULL, '{\"en-US\": \"Supplier\", \"zh-CN\": \"供应商\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (36, 'product.list.unitCost', 'product', 'list', 'active', NULL, '{\"en-US\": \"Unit Cost\", \"zh-CN\": \"单位成本\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (37, 'product.list.productImage', 'product', 'list', 'active', NULL, '{\"en-US\": \"Product Image\", \"zh-CN\": \"产品图片\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (38, 'product.list.unitCostUsd', 'product', 'list', 'active', NULL, '{\"en-US\": \"Unit Cost (USD)\", \"zh-CN\": \"单位成本(USD)\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (39, 'product.list.fnsku', 'product', 'list', 'active', NULL, '{\"en-US\": \"FNSKU\", \"zh-CN\": \"FNSKU\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (40, 'product.list.remark', 'product', 'list', 'active', NULL, '{\"en-US\": \"Remark\", \"zh-CN\": \"备注\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (41, 'product.list.skuDetails', 'product', 'list', 'active', NULL, '{\"en-US\": \"SKU Details\", \"zh-CN\": \"SKU详情\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (42, 'product.list.createdAt', 'product', 'list', 'active', NULL, '{\"en-US\": \"Created At\", \"zh-CN\": \"创建时间\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (43, 'product.list.updatedAt', 'product', 'list', 'active', NULL, '{\"en-US\": \"Updated At\", \"zh-CN\": \"更新时间\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (44, 'product.list.auditLogs', 'product', 'list', 'active', NULL, '{\"en-US\": \"Audit Logs\", \"zh-CN\": \"审计日志\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (45, 'product.list.time', 'product', 'list', 'active', NULL, '{\"en-US\": \"Time\", \"zh-CN\": \"时间\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (46, 'product.list.action', 'product', 'list', 'active', NULL, '{\"en-US\": \"Action\", \"zh-CN\": \"动作\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (47, 'product.list.changes', 'product', 'list', 'active', NULL, '{\"en-US\": \"Changes\", \"zh-CN\": \"变更\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (48, 'product.list.children', 'product', 'list', 'active', NULL, '{\"en-US\": \"Children\", \"zh-CN\": \"子产品\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (49, 'product.list.noChildren', 'product', 'list', 'active', NULL, '{\"en-US\": \"No child products\", \"zh-CN\": \"暂无子产品\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (50, 'product.list.productInfo', 'product', 'list', 'active', NULL, '{\"en-US\": \"Product Info\", \"zh-CN\": \"产品信息\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (51, 'product.list.supplierInfo', 'product', 'list', 'active', NULL, '{\"en-US\": \"Supplier Info\", \"zh-CN\": \"供应商信息\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (52, 'product.list.inventoryInfo', 'product', 'list', 'active', NULL, '{\"en-US\": \"Inventory Info\", \"zh-CN\": \"库存信息\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (53, 'product.list.priceInfo', 'product', 'list', 'active', NULL, '{\"en-US\": \"Price & Cost\", \"zh-CN\": \"价格/成本\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (54, 'product.list.statusInfo', 'product', 'list', 'active', NULL, '{\"en-US\": \"Status & Site\", \"zh-CN\": \"状态/站点\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (55, 'product.list.createdInfo', 'product', 'list', 'active', NULL, '{\"en-US\": \"Created Info\", \"zh-CN\": \"创建信息\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (56, 'product.list.available', 'product', 'list', 'active', NULL, '{\"en-US\": \"Available\", \"zh-CN\": \"可用\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (57, 'product.list.reserved', 'product', 'list', 'active', NULL, '{\"en-US\": \"Reserved\", \"zh-CN\": \"占用\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (58, 'product.list.inbound', 'product', 'list', 'active', NULL, '{\"en-US\": \"Inbound\", \"zh-CN\": \"在途\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (59, 'product.list.weight', 'product', 'list', 'active', NULL, '{\"en-US\": \"Weight\", \"zh-CN\": \"重量\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (60, 'product.list.dimensions', 'product', 'list', 'active', NULL, '{\"en-US\": \"Dimensions\", \"zh-CN\": \"尺寸\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (61, 'product.list.createdBy', 'product', 'list', 'active', NULL, '{\"en-US\": \"Created By\", \"zh-CN\": \"创建人\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (62, 'product.list.comboTag', 'product', 'list', 'active', NULL, '{\"en-US\": \"COMBO\", \"zh-CN\": \"组合\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (63, 'product.list.deleteConfirm', 'product', 'list', 'active', NULL, '{\"en-US\": \"Are you sure to delete SKU \\\"{sku}\\\"?\", \"zh-CN\": \"确认删除SKU \\\"{sku}\\\"？\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (64, 'product.list.createdSuccess', 'product', 'list', 'active', NULL, '{\"en-US\": \"SKU created successfully\", \"zh-CN\": \"SKU创建成功\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (65, 'product.list.updatedSuccess', 'product', 'list', 'active', NULL, '{\"en-US\": \"SKU updated successfully\", \"zh-CN\": \"SKU更新成功\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (66, 'product.list.deletedSuccess', 'product', 'list', 'active', NULL, '{\"en-US\": \"SKU deleted successfully\", \"zh-CN\": \"SKU删除成功\"}', '2026-01-19 15:26:10', '2026-01-19 15:26:10');
INSERT INTO `field_label` VALUES (67, 'global.submit', NULL, NULL, 'active', NULL, '{\"en-US\": \"Submit\", \"zh-CN\": \"提交\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (68, 'global.open', NULL, NULL, 'active', NULL, '{\"en-US\": \"Open\", \"zh-CN\": \"打开\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (69, 'global.close', NULL, NULL, 'active', NULL, '{\"en-US\": \"Close\", \"zh-CN\": \"关闭\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (70, 'global.enable', NULL, NULL, 'active', NULL, '{\"en-US\": \"Enable\", \"zh-CN\": \"启用\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (71, 'global.disable', NULL, NULL, 'active', NULL, '{\"en-US\": \"Disable\", \"zh-CN\": \"禁用\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (72, 'global.back', NULL, NULL, 'active', NULL, '{\"en-US\": \"Back\", \"zh-CN\": \"返回\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (73, 'global.next', NULL, NULL, 'active', NULL, '{\"en-US\": \"Next\", \"zh-CN\": \"下一步\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (74, 'global.previous', NULL, NULL, 'active', NULL, '{\"en-US\": \"Previous\", \"zh-CN\": \"上一步\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (75, 'global.ok', NULL, NULL, 'active', NULL, '{\"en-US\": \"OK\", \"zh-CN\": \"确定\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (76, 'global.yes', NULL, NULL, 'active', NULL, '{\"en-US\": \"Yes\", \"zh-CN\": \"是\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (77, 'global.no', NULL, NULL, 'active', NULL, '{\"en-US\": \"No\", \"zh-CN\": \"否\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (78, 'global.clear', NULL, NULL, 'active', NULL, '{\"en-US\": \"Clear\", \"zh-CN\": \"清空\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (79, 'global.filter', NULL, NULL, 'active', NULL, '{\"en-US\": \"Filter\", \"zh-CN\": \"筛选\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (80, 'global.download', NULL, NULL, 'active', NULL, '{\"en-US\": \"Download\", \"zh-CN\": \"下载\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (81, 'global.upload', NULL, NULL, 'active', NULL, '{\"en-US\": \"Upload\", \"zh-CN\": \"上传\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (82, 'global.copy', NULL, NULL, 'active', NULL, '{\"en-US\": \"Copy\", \"zh-CN\": \"复制\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (83, 'global.create', NULL, NULL, 'active', NULL, '{\"en-US\": \"Create\", \"zh-CN\": \"创建\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (84, 'global.update', NULL, NULL, 'active', NULL, '{\"en-US\": \"Update\", \"zh-CN\": \"更新\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (85, 'global.deleteConfirm', NULL, NULL, 'active', NULL, '{\"en-US\": \"Confirm delete?\", \"zh-CN\": \"确认删除？\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (86, 'global.success', NULL, NULL, 'active', NULL, '{\"en-US\": \"Success\", \"zh-CN\": \"操作成功\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (87, 'global.failed', NULL, NULL, 'active', NULL, '{\"en-US\": \"Failed\", \"zh-CN\": \"操作失败\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (88, 'global.loading', NULL, NULL, 'active', NULL, '{\"en-US\": \"Loading\", \"zh-CN\": \"加载中\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (89, 'global.empty', NULL, NULL, 'active', NULL, '{\"en-US\": \"No data\", \"zh-CN\": \"暂无数据\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (90, 'global.retry', NULL, NULL, 'active', NULL, '{\"en-US\": \"Retry\", \"zh-CN\": \"重试\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (91, 'global.confirmSubmit', NULL, NULL, 'active', NULL, '{\"en-US\": \"Confirm submit?\", \"zh-CN\": \"确认提交？\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (92, 'global.confirmCancel', NULL, NULL, 'active', NULL, '{\"en-US\": \"Confirm cancel?\", \"zh-CN\": \"确认取消？\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (93, 'global.confirmClose', NULL, NULL, 'active', NULL, '{\"en-US\": \"Confirm close?\", \"zh-CN\": \"确认关闭？\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (94, 'global.confirmUpdate', NULL, NULL, 'active', NULL, '{\"en-US\": \"Confirm update?\", \"zh-CN\": \"确认更新？\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (95, 'global.confirmSave', NULL, NULL, 'active', NULL, '{\"en-US\": \"Confirm save?\", \"zh-CN\": \"确认保存？\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (96, 'global.warningTitle', NULL, NULL, 'active', NULL, '{\"en-US\": \"Warning\", \"zh-CN\": \"提示\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (97, 'global.errorTitle', NULL, NULL, 'active', NULL, '{\"en-US\": \"Error\", \"zh-CN\": \"错误\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (98, 'global.infoTitle', NULL, NULL, 'active', NULL, '{\"en-US\": \"Info\", \"zh-CN\": \"信息\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (99, 'global.total', NULL, NULL, 'active', NULL, '{\"en-US\": \"Total\", \"zh-CN\": \"总计\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (100, 'global.page', NULL, NULL, 'active', NULL, '{\"en-US\": \"Page\", \"zh-CN\": \"页\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (101, 'global.pageSize', NULL, NULL, 'active', NULL, '{\"en-US\": \"Page Size\", \"zh-CN\": \"每页\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (102, 'global.jumpTo', NULL, NULL, 'active', NULL, '{\"en-US\": \"Jump to\", \"zh-CN\": \"跳转\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (103, 'global.firstPage', NULL, NULL, 'active', NULL, '{\"en-US\": \"First\", \"zh-CN\": \"首页\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (104, 'global.lastPage', NULL, NULL, 'active', NULL, '{\"en-US\": \"Last\", \"zh-CN\": \"末页\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (105, 'global.from', NULL, NULL, 'active', NULL, '{\"en-US\": \"From\", \"zh-CN\": \"从\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (106, 'global.to', NULL, NULL, 'active', NULL, '{\"en-US\": \"To\", \"zh-CN\": \"到\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (107, 'global.startDate', NULL, NULL, 'active', NULL, '{\"en-US\": \"Start Date\", \"zh-CN\": \"开始日期\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (108, 'global.endDate', NULL, NULL, 'active', NULL, '{\"en-US\": \"End Date\", \"zh-CN\": \"结束日期\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (109, 'global.dateRange', NULL, NULL, 'active', NULL, '{\"en-US\": \"Date Range\", \"zh-CN\": \"日期范围\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (110, 'global.exportSuccess', NULL, NULL, 'active', NULL, '{\"en-US\": \"Exported successfully\", \"zh-CN\": \"导出成功\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (111, 'global.importSuccess', NULL, NULL, 'active', NULL, '{\"en-US\": \"Imported successfully\", \"zh-CN\": \"导入成功\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (112, 'global.importFailed', NULL, NULL, 'active', NULL, '{\"en-US\": \"Import failed\", \"zh-CN\": \"导入失败\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (113, 'global.login', NULL, NULL, 'active', NULL, '{\"en-US\": \"Login\", \"zh-CN\": \"登录\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (114, 'global.logout', NULL, NULL, 'active', NULL, '{\"en-US\": \"Logout\", \"zh-CN\": \"退出登录\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (115, 'global.username', NULL, NULL, 'active', NULL, '{\"en-US\": \"Username\", \"zh-CN\": \"用户名\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (116, 'global.password', NULL, NULL, 'active', NULL, '{\"en-US\": \"Password\", \"zh-CN\": \"密码\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (117, 'global.rememberMe', NULL, NULL, 'active', NULL, '{\"en-US\": \"Remember me\", \"zh-CN\": \"记住我\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (118, 'global.id', NULL, NULL, 'active', NULL, '{\"en-US\": \"ID\", \"zh-CN\": \"ID\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (119, 'global.name', NULL, NULL, 'active', NULL, '{\"en-US\": \"Name\", \"zh-CN\": \"名称\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (120, 'global.code', NULL, NULL, 'active', NULL, '{\"en-US\": \"Code\", \"zh-CN\": \"编码\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (121, 'global.type', NULL, NULL, 'active', NULL, '{\"en-US\": \"Type\", \"zh-CN\": \"类型\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (122, 'global.status', NULL, NULL, 'active', NULL, '{\"en-US\": \"Status\", \"zh-CN\": \"状态\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (123, 'global.active', NULL, NULL, 'active', NULL, '{\"en-US\": \"Active\", \"zh-CN\": \"启用\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (124, 'global.inactive', NULL, NULL, 'active', NULL, '{\"en-US\": \"Inactive\", \"zh-CN\": \"停用\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (125, 'global.enabled', NULL, NULL, 'active', NULL, '{\"en-US\": \"Enabled\", \"zh-CN\": \"启用\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (126, 'global.disabled', NULL, NULL, 'active', NULL, '{\"en-US\": \"Disabled\", \"zh-CN\": \"禁用\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (127, 'global.remark', NULL, NULL, 'active', NULL, '{\"en-US\": \"Remark\", \"zh-CN\": \"备注\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (128, 'global.description', NULL, NULL, 'active', NULL, '{\"en-US\": \"Description\", \"zh-CN\": \"说明\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (129, 'global.createdAt', NULL, NULL, 'active', NULL, '{\"en-US\": \"Created At\", \"zh-CN\": \"创建时间\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (130, 'global.updatedAt', NULL, NULL, 'active', NULL, '{\"en-US\": \"Updated At\", \"zh-CN\": \"更新时间\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (131, 'global.createdBy', NULL, NULL, 'active', NULL, '{\"en-US\": \"Created By\", \"zh-CN\": \"创建人\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (132, 'global.updatedBy', NULL, NULL, 'active', NULL, '{\"en-US\": \"Updated By\", \"zh-CN\": \"更新人\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (133, 'global.start', NULL, NULL, 'active', NULL, '{\"en-US\": \"Start\", \"zh-CN\": \"开始\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (134, 'global.end', NULL, NULL, 'active', NULL, '{\"en-US\": \"End\", \"zh-CN\": \"结束\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (135, 'global.startTime', NULL, NULL, 'active', NULL, '{\"en-US\": \"Start Time\", \"zh-CN\": \"开始时间\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (136, 'global.endTime', NULL, NULL, 'active', NULL, '{\"en-US\": \"End Time\", \"zh-CN\": \"结束时间\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (137, 'global.date', NULL, NULL, 'active', NULL, '{\"en-US\": \"Date\", \"zh-CN\": \"日期\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (138, 'global.time', NULL, NULL, 'active', NULL, '{\"en-US\": \"Time\", \"zh-CN\": \"时间\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (139, 'global.amount', NULL, NULL, 'active', NULL, '{\"en-US\": \"Amount\", \"zh-CN\": \"金额\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (140, 'global.totalAmount', NULL, NULL, 'active', NULL, '{\"en-US\": \"Total Amount\", \"zh-CN\": \"总金额\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (141, 'global.quantity', NULL, NULL, 'active', NULL, '{\"en-US\": \"Quantity\", \"zh-CN\": \"数量\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (142, 'global.unit', NULL, NULL, 'active', NULL, '{\"en-US\": \"Unit\", \"zh-CN\": \"单位\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (143, 'global.unitPrice', NULL, NULL, 'active', NULL, '{\"en-US\": \"Unit Price\", \"zh-CN\": \"单价\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (144, 'global.unitCost', NULL, NULL, 'active', NULL, '{\"en-US\": \"Unit Cost\", \"zh-CN\": \"单位成本\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (145, 'global.totalCost', NULL, NULL, 'active', NULL, '{\"en-US\": \"Total Cost\", \"zh-CN\": \"总成本\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (146, 'global.totalPrice', NULL, NULL, 'active', NULL, '{\"en-US\": \"Total Price\", \"zh-CN\": \"总价\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (147, 'global.currency', NULL, NULL, 'active', NULL, '{\"en-US\": \"Currency\", \"zh-CN\": \"币种\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (148, 'global.weight', NULL, NULL, 'active', NULL, '{\"en-US\": \"Weight\", \"zh-CN\": \"重量\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (149, 'global.dimensions', NULL, NULL, 'active', NULL, '{\"en-US\": \"Dimensions\", \"zh-CN\": \"尺寸\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (150, 'global.all', NULL, NULL, 'active', NULL, '{\"en-US\": \"All\", \"zh-CN\": \"全部\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (151, 'global.none', NULL, NULL, 'active', NULL, '{\"en-US\": \"None\", \"zh-CN\": \"无\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (152, 'global.select', NULL, NULL, 'active', NULL, '{\"en-US\": \"Select\", \"zh-CN\": \"选择\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (153, 'global.selectAll', NULL, NULL, 'active', NULL, '{\"en-US\": \"Select All\", \"zh-CN\": \"全选\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (154, 'global.clearSelection', NULL, NULL, 'active', NULL, '{\"en-US\": \"Clear Selection\", \"zh-CN\": \"清除选择\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (155, 'global.addItem', NULL, NULL, 'active', NULL, '{\"en-US\": \"Add Item\", \"zh-CN\": \"添加项\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (156, 'global.removeItem', NULL, NULL, 'active', NULL, '{\"en-US\": \"Remove Item\", \"zh-CN\": \"移除项\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (157, 'global.uploadFile', NULL, NULL, 'active', NULL, '{\"en-US\": \"Upload File\", \"zh-CN\": \"上传文件\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (158, 'global.chooseFile', NULL, NULL, 'active', NULL, '{\"en-US\": \"Choose File\", \"zh-CN\": \"选择文件\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (159, 'global.downloadTemplate', NULL, NULL, 'active', NULL, '{\"en-US\": \"Download Template\", \"zh-CN\": \"下载模板\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (160, 'global.uploadSuccess', NULL, NULL, 'active', NULL, '{\"en-US\": \"Upload successful\", \"zh-CN\": \"上传成功\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (161, 'global.uploadFailed', NULL, NULL, 'active', NULL, '{\"en-US\": \"Upload failed\", \"zh-CN\": \"上传失败\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (162, 'global.saveSuccess', NULL, NULL, 'active', NULL, '{\"en-US\": \"Saved successfully\", \"zh-CN\": \"保存成功\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (163, 'global.updateSuccess', NULL, NULL, 'active', NULL, '{\"en-US\": \"Updated successfully\", \"zh-CN\": \"更新成功\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (164, 'global.createSuccess', NULL, NULL, 'active', NULL, '{\"en-US\": \"Created successfully\", \"zh-CN\": \"创建成功\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (165, 'global.deleteSuccess', NULL, NULL, 'active', NULL, '{\"en-US\": \"Deleted successfully\", \"zh-CN\": \"删除成功\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (166, 'global.confirmDeleteWithName', NULL, NULL, 'active', NULL, '{\"en-US\": \"Confirm delete \\\"{name}\\\"?\", \"zh-CN\": \"确认删除“{name}”？\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (167, 'global.confirmSubmitWithName', NULL, NULL, 'active', NULL, '{\"en-US\": \"Confirm submit \\\"{name}\\\"?\", \"zh-CN\": \"确认提交“{name}”？\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (168, 'global.supplier', NULL, NULL, 'active', NULL, '{\"en-US\": \"Supplier\", \"zh-CN\": \"供应商\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (169, 'global.customer', NULL, NULL, 'active', NULL, '{\"en-US\": \"Customer\", \"zh-CN\": \"客户\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (170, 'global.warehouse', NULL, NULL, 'active', NULL, '{\"en-US\": \"Warehouse\", \"zh-CN\": \"仓库\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (171, 'global.order', NULL, NULL, 'active', NULL, '{\"en-US\": \"Order\", \"zh-CN\": \"订单\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (172, 'global.orderNo', NULL, NULL, 'active', NULL, '{\"en-US\": \"Order No.\", \"zh-CN\": \"订单号\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (173, 'global.purchaseOrder', NULL, NULL, 'active', NULL, '{\"en-US\": \"Purchase Order\", \"zh-CN\": \"采购单\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (174, 'global.salesOrder', NULL, NULL, 'active', NULL, '{\"en-US\": \"Sales Order\", \"zh-CN\": \"销售单\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (175, 'global.shipment', NULL, NULL, 'active', NULL, '{\"en-US\": \"Shipment\", \"zh-CN\": \"货件\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (176, 'global.trackingNo', NULL, NULL, 'active', NULL, '{\"en-US\": \"Tracking No.\", \"zh-CN\": \"物流单号\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (177, 'global.invoice', NULL, NULL, 'active', NULL, '{\"en-US\": \"Invoice\", \"zh-CN\": \"发票\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (178, 'global.payment', NULL, NULL, 'active', NULL, '{\"en-US\": \"Payment\", \"zh-CN\": \"付款\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (179, 'global.receipt', NULL, NULL, 'active', NULL, '{\"en-US\": \"Receipt\", \"zh-CN\": \"收据\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (180, 'global.balance', NULL, NULL, 'active', NULL, '{\"en-US\": \"Balance\", \"zh-CN\": \"余额\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (181, 'global.cost', NULL, NULL, 'active', NULL, '{\"en-US\": \"Cost\", \"zh-CN\": \"成本\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (182, 'global.profit', NULL, NULL, 'active', NULL, '{\"en-US\": \"Profit\", \"zh-CN\": \"利润\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (183, 'global.tax', NULL, NULL, 'active', NULL, '{\"en-US\": \"Tax\", \"zh-CN\": \"税费\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (184, 'global.discount', NULL, NULL, 'active', NULL, '{\"en-US\": \"Discount\", \"zh-CN\": \"折扣\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (185, 'global.fee', NULL, NULL, 'active', NULL, '{\"en-US\": \"Fee\", \"zh-CN\": \"费用\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (186, 'global.freight', NULL, NULL, 'active', NULL, '{\"en-US\": \"Freight\", \"zh-CN\": \"运费\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (187, 'global.origin', NULL, NULL, 'active', NULL, '{\"en-US\": \"Origin\", \"zh-CN\": \"起点\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (188, 'global.destination', NULL, NULL, 'active', NULL, '{\"en-US\": \"Destination\", \"zh-CN\": \"目的地\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (189, 'global.etd', NULL, NULL, 'active', NULL, '{\"en-US\": \"ETD\", \"zh-CN\": \"预计发货\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (190, 'global.eta', NULL, NULL, 'active', NULL, '{\"en-US\": \"ETA\", \"zh-CN\": \"预计到达\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (191, 'global.draft', NULL, NULL, 'active', NULL, '{\"en-US\": \"Draft\", \"zh-CN\": \"草稿\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (192, 'global.pending', NULL, NULL, 'active', NULL, '{\"en-US\": \"Pending\", \"zh-CN\": \"待处理\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (193, 'global.processing', NULL, NULL, 'active', NULL, '{\"en-US\": \"Processing\", \"zh-CN\": \"处理中\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (194, 'global.approved', NULL, NULL, 'active', NULL, '{\"en-US\": \"Approved\", \"zh-CN\": \"已审核\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (195, 'global.rejected', NULL, NULL, 'active', NULL, '{\"en-US\": \"Rejected\", \"zh-CN\": \"已拒绝\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (196, 'global.completed', NULL, NULL, 'active', NULL, '{\"en-US\": \"Completed\", \"zh-CN\": \"已完成\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (197, 'global.cancelled', NULL, NULL, 'active', NULL, '{\"en-US\": \"Cancelled\", \"zh-CN\": \"已取消\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (198, 'global.closed', NULL, NULL, 'active', NULL, '{\"en-US\": \"Closed\", \"zh-CN\": \"已关闭\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (199, 'global.inProgress', NULL, NULL, 'active', NULL, '{\"en-US\": \"In Progress\", \"zh-CN\": \"进行中\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (200, 'global.onHold', NULL, NULL, 'active', NULL, '{\"en-US\": \"On Hold\", \"zh-CN\": \"暂停\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (201, 'global.partial', NULL, NULL, 'active', NULL, '{\"en-US\": \"Partial\", \"zh-CN\": \"部分\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (202, 'global.overdue', NULL, NULL, 'active', NULL, '{\"en-US\": \"Overdue\", \"zh-CN\": \"逾期\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (203, 'global.delivered', NULL, NULL, 'active', NULL, '{\"en-US\": \"Delivered\", \"zh-CN\": \"已送达\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (204, 'global.received', NULL, NULL, 'active', NULL, '{\"en-US\": \"Received\", \"zh-CN\": \"已收货\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (205, 'global.shipped', NULL, NULL, 'active', NULL, '{\"en-US\": \"Shipped\", \"zh-CN\": \"已发货\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (206, 'global.orderPlaced', NULL, NULL, 'active', NULL, '{\"en-US\": \"Ordered\", \"zh-CN\": \"已下单\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (207, 'global.refunded', NULL, NULL, 'active', NULL, '{\"en-US\": \"Refunded\", \"zh-CN\": \"已退款\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (208, 'global.returned', NULL, NULL, 'active', NULL, '{\"en-US\": \"Returned\", \"zh-CN\": \"已退货\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (209, 'global.inbound', NULL, NULL, 'active', NULL, '{\"en-US\": \"Inbound\", \"zh-CN\": \"入库\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (210, 'global.outbound', NULL, NULL, 'active', NULL, '{\"en-US\": \"Outbound\", \"zh-CN\": \"出库\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (211, 'global.adjustment', NULL, NULL, 'active', NULL, '{\"en-US\": \"Adjustment\", \"zh-CN\": \"调整\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (212, 'global.transfer', NULL, NULL, 'active', NULL, '{\"en-US\": \"Transfer\", \"zh-CN\": \"调拨\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (213, 'global.stocktake', NULL, NULL, 'active', NULL, '{\"en-US\": \"Stock Take\", \"zh-CN\": \"盘点\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (214, 'global.damage', NULL, NULL, 'active', NULL, '{\"en-US\": \"Damage\", \"zh-CN\": \"报损\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (215, 'global.purchaseReceipt', NULL, NULL, 'active', NULL, '{\"en-US\": \"Purchase Receipt\", \"zh-CN\": \"采购入库\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (216, 'global.salesShipment', NULL, NULL, 'active', NULL, '{\"en-US\": \"Sales Shipment\", \"zh-CN\": \"销售出库\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (217, 'global.transferIn', NULL, NULL, 'active', NULL, '{\"en-US\": \"Transfer In\", \"zh-CN\": \"调拨入库\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (218, 'global.transferOut', NULL, NULL, 'active', NULL, '{\"en-US\": \"Transfer Out\", \"zh-CN\": \"调拨出库\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (219, 'global.role', NULL, NULL, 'active', NULL, '{\"en-US\": \"Role\", \"zh-CN\": \"角色\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (220, 'global.permission', NULL, NULL, 'active', NULL, '{\"en-US\": \"Permission\", \"zh-CN\": \"权限\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (221, 'global.user', NULL, NULL, 'active', NULL, '{\"en-US\": \"User\", \"zh-CN\": \"用户\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (222, 'global.email', NULL, NULL, 'active', NULL, '{\"en-US\": \"Email\", \"zh-CN\": \"邮箱\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (223, 'global.phone', NULL, NULL, 'active', NULL, '{\"en-US\": \"Phone\", \"zh-CN\": \"电话\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (224, 'global.address', NULL, NULL, 'active', NULL, '{\"en-US\": \"Address\", \"zh-CN\": \"地址\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (225, 'global.contact', NULL, NULL, 'active', NULL, '{\"en-US\": \"Contact\", \"zh-CN\": \"联系人\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (226, 'global.site', NULL, NULL, 'active', NULL, '{\"en-US\": \"Site\", \"zh-CN\": \"站点\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (227, 'global.channel', NULL, NULL, 'active', NULL, '{\"en-US\": \"Channel\", \"zh-CN\": \"渠道\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (228, 'global.marketplace', NULL, NULL, 'active', NULL, '{\"en-US\": \"Marketplace\", \"zh-CN\": \"市场\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (229, 'global.platform', NULL, NULL, 'active', NULL, '{\"en-US\": \"Platform\", \"zh-CN\": \"平台\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (230, 'global.region', NULL, NULL, 'active', NULL, '{\"en-US\": \"Region\", \"zh-CN\": \"区域\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (231, 'global.country', NULL, NULL, 'active', NULL, '{\"en-US\": \"Country\", \"zh-CN\": \"国家\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (232, 'global.currencyCNY', NULL, NULL, 'active', NULL, '{\"en-US\": \"CNY\", \"zh-CN\": \"人民币\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (233, 'global.currencyUSD', NULL, NULL, 'active', NULL, '{\"en-US\": \"USD\", \"zh-CN\": \"美元\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (234, 'global.currencyEUR', NULL, NULL, 'active', NULL, '{\"en-US\": \"EUR\", \"zh-CN\": \"欧元\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (235, 'global.currencyGBP', NULL, NULL, 'active', NULL, '{\"en-US\": \"GBP\", \"zh-CN\": \"英镑\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (236, 'global.currencyJPY', NULL, NULL, 'active', NULL, '{\"en-US\": \"JPY\", \"zh-CN\": \"日元\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (237, 'global.currencyAUD', NULL, NULL, 'active', NULL, '{\"en-US\": \"AUD\", \"zh-CN\": \"澳元\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (238, 'global.currencyCAD', NULL, NULL, 'active', NULL, '{\"en-US\": \"CAD\", \"zh-CN\": \"加元\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (239, 'global.taxRate', NULL, NULL, 'active', NULL, '{\"en-US\": \"Tax Rate\", \"zh-CN\": \"税率\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (240, 'global.taxIncluded', NULL, NULL, 'active', NULL, '{\"en-US\": \"Tax Included\", \"zh-CN\": \"含税\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (241, 'global.taxExcluded', NULL, NULL, 'active', NULL, '{\"en-US\": \"Tax Excluded\", \"zh-CN\": \"不含税\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (242, 'global.paymentMethod', NULL, NULL, 'active', NULL, '{\"en-US\": \"Payment Method\", \"zh-CN\": \"支付方式\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (243, 'global.shippingMethod', NULL, NULL, 'active', NULL, '{\"en-US\": \"Shipping Method\", \"zh-CN\": \"运输方式\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (244, 'global.shippingCarrier', NULL, NULL, 'active', NULL, '{\"en-US\": \"Carrier\", \"zh-CN\": \"承运商\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (245, 'global.air', NULL, NULL, 'active', NULL, '{\"en-US\": \"Air\", \"zh-CN\": \"空运\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (246, 'global.sea', NULL, NULL, 'active', NULL, '{\"en-US\": \"Sea\", \"zh-CN\": \"海运\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (247, 'global.land', NULL, NULL, 'active', NULL, '{\"en-US\": \"Land\", \"zh-CN\": \"陆运\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (248, 'global.express', NULL, NULL, 'active', NULL, '{\"en-US\": \"Express\", \"zh-CN\": \"快递\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (249, 'global.standard', NULL, NULL, 'active', NULL, '{\"en-US\": \"Standard\", \"zh-CN\": \"标准\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (250, 'global.priority', NULL, NULL, 'active', NULL, '{\"en-US\": \"Priority\", \"zh-CN\": \"优先\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (251, 'global.fba', NULL, NULL, 'active', NULL, '{\"en-US\": \"FBA\", \"zh-CN\": \"FBA\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (252, 'global.thirdParty', NULL, NULL, 'active', NULL, '{\"en-US\": \"Third Party\", \"zh-CN\": \"第三方仓\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (253, 'global.ownWarehouse', NULL, NULL, 'active', NULL, '{\"en-US\": \"Own Warehouse\", \"zh-CN\": \"自有仓\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (254, 'global.availableQty', NULL, NULL, 'active', NULL, '{\"en-US\": \"Available Qty\", \"zh-CN\": \"可用库存\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (255, 'global.reservedQty', NULL, NULL, 'active', NULL, '{\"en-US\": \"Reserved Qty\", \"zh-CN\": \"占用库存\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (256, 'global.inboundQty', NULL, NULL, 'active', NULL, '{\"en-US\": \"Inbound Qty\", \"zh-CN\": \"在途库存\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (257, 'global.safetyStock', NULL, NULL, 'active', NULL, '{\"en-US\": \"Safety Stock\", \"zh-CN\": \"安全库存\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (258, 'global.reorderPoint', NULL, NULL, 'active', NULL, '{\"en-US\": \"Reorder Point\", \"zh-CN\": \"补货点\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (259, 'global.leadTime', NULL, NULL, 'active', NULL, '{\"en-US\": \"Lead Time\", \"zh-CN\": \"采购周期\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (260, 'global.moQ', NULL, NULL, 'active', NULL, '{\"en-US\": \"MOQ\", \"zh-CN\": \"最小起订量\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (261, 'global.sku', NULL, NULL, 'active', NULL, '{\"en-US\": \"SKU\", \"zh-CN\": \"SKU\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (262, 'global.asin', NULL, NULL, 'active', NULL, '{\"en-US\": \"ASIN\", \"zh-CN\": \"ASIN\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (263, 'global.fnsku', NULL, NULL, 'active', NULL, '{\"en-US\": \"FNSKU\", \"zh-CN\": \"FNSKU\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (264, 'global.brand', NULL, NULL, 'active', NULL, '{\"en-US\": \"Brand\", \"zh-CN\": \"品牌\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (265, 'global.category', NULL, NULL, 'active', NULL, '{\"en-US\": \"Category\", \"zh-CN\": \"分类\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (266, 'global.batchNo', NULL, NULL, 'active', NULL, '{\"en-US\": \"Batch No.\", \"zh-CN\": \"批次号\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (267, 'global.lotNo', NULL, NULL, 'active', NULL, '{\"en-US\": \"Lot No.\", \"zh-CN\": \"批号\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (268, 'global.serialNo', NULL, NULL, 'active', NULL, '{\"en-US\": \"Serial No.\", \"zh-CN\": \"序列号\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (269, 'global.barcode', NULL, NULL, 'active', NULL, '{\"en-US\": \"Barcode\", \"zh-CN\": \"条码\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (270, 'global.expiryDate', NULL, NULL, 'active', NULL, '{\"en-US\": \"Expiry Date\", \"zh-CN\": \"有效期\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (271, 'global.manufactureDate', NULL, NULL, 'active', NULL, '{\"en-US\": \"Manufacture Date\", \"zh-CN\": \"生产日期\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (272, 'global.packing', NULL, NULL, 'active', NULL, '{\"en-US\": \"Packaging\", \"zh-CN\": \"包装\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (273, 'global.material', NULL, NULL, 'active', NULL, '{\"en-US\": \"Material\", \"zh-CN\": \"材质\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (274, 'global.color', NULL, NULL, 'active', NULL, '{\"en-US\": \"Color\", \"zh-CN\": \"颜色\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (275, 'global.size', NULL, NULL, 'active', NULL, '{\"en-US\": \"Size\", \"zh-CN\": \"尺寸\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (276, 'global.volume', NULL, NULL, 'active', NULL, '{\"en-US\": \"Volume\", \"zh-CN\": \"体积\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (277, 'global.cbm', NULL, NULL, 'active', NULL, '{\"en-US\": \"CBM\", \"zh-CN\": \"方数\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (278, 'global.container', NULL, NULL, 'active', NULL, '{\"en-US\": \"Container\", \"zh-CN\": \"集装箱\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (279, 'global.pallet', NULL, NULL, 'active', NULL, '{\"en-US\": \"Pallet\", \"zh-CN\": \"托盘\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (280, 'global.carton', NULL, NULL, 'active', NULL, '{\"en-US\": \"Carton\", \"zh-CN\": \"箱\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (281, 'global.packagingMaterial', NULL, NULL, 'active', NULL, '{\"en-US\": \"Packaging Material\", \"zh-CN\": \"包材\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (282, 'global.consumed', NULL, NULL, 'active', NULL, '{\"en-US\": \"Consumed\", \"zh-CN\": \"已消耗\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (283, 'global.remaining', NULL, NULL, 'active', NULL, '{\"en-US\": \"Remaining\", \"zh-CN\": \"剩余\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (284, 'global.increase', NULL, NULL, 'active', NULL, '{\"en-US\": \"Increase\", \"zh-CN\": \"增加\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (285, 'global.decrease', NULL, NULL, 'active', NULL, '{\"en-US\": \"Decrease\", \"zh-CN\": \"减少\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (286, 'global.manual', NULL, NULL, 'active', NULL, '{\"en-US\": \"Manual\", \"zh-CN\": \"手动\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (287, 'global.system', NULL, NULL, 'active', NULL, '{\"en-US\": \"System\", \"zh-CN\": \"系统\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (288, 'global.reason', NULL, NULL, 'active', NULL, '{\"en-US\": \"Reason\", \"zh-CN\": \"原因\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (289, 'global.referenceNo', NULL, NULL, 'active', NULL, '{\"en-US\": \"Reference No.\", \"zh-CN\": \"参考单号\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (290, 'global.referenceType', NULL, NULL, 'active', NULL, '{\"en-US\": \"Reference Type\", \"zh-CN\": \"参考类型\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (291, 'global.notes', NULL, NULL, 'active', NULL, '{\"en-US\": \"Notes\", \"zh-CN\": \"备注说明\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (292, 'global.balanceQty', NULL, NULL, 'active', NULL, '{\"en-US\": \"Balance Qty\", \"zh-CN\": \"结存数量\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (293, 'global.actualQty', NULL, NULL, 'active', NULL, '{\"en-US\": \"Actual Qty\", \"zh-CN\": \"实际数量\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (294, 'global.expectedQty', NULL, NULL, 'active', NULL, '{\"en-US\": \"Expected Qty\", \"zh-CN\": \"预计数量\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (295, 'global.difference', NULL, NULL, 'active', NULL, '{\"en-US\": \"Difference\", \"zh-CN\": \"差异\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (296, 'global.pick', NULL, NULL, 'active', NULL, '{\"en-US\": \"Pick\", \"zh-CN\": \"拣货\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (297, 'global.pack', NULL, NULL, 'active', NULL, '{\"en-US\": \"Pack\", \"zh-CN\": \"打包\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (298, 'global.ship', NULL, NULL, 'active', NULL, '{\"en-US\": \"Ship\", \"zh-CN\": \"发货\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (299, 'global.receive', NULL, NULL, 'active', NULL, '{\"en-US\": \"Receive\", \"zh-CN\": \"收货\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (300, 'global.payable', NULL, NULL, 'active', NULL, '{\"en-US\": \"Payable\", \"zh-CN\": \"应付\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (301, 'global.receivable', NULL, NULL, 'active', NULL, '{\"en-US\": \"Receivable\", \"zh-CN\": \"应收\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (302, 'global.settlement', NULL, NULL, 'active', NULL, '{\"en-US\": \"Settlement\", \"zh-CN\": \"结算\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (303, 'global.balanceSheet', NULL, NULL, 'active', NULL, '{\"en-US\": \"Balance Sheet\", \"zh-CN\": \"资产负债表\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (304, 'global.profitLoss', NULL, NULL, 'active', NULL, '{\"en-US\": \"Profit & Loss\", \"zh-CN\": \"损益表\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (305, 'global.cashFlow', NULL, NULL, 'active', NULL, '{\"en-US\": \"Cash Flow\", \"zh-CN\": \"现金流\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (306, 'global.inventory', NULL, NULL, 'active', NULL, '{\"en-US\": \"Inventory\", \"zh-CN\": \"库存\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (307, 'global.procurement', NULL, NULL, 'active', NULL, '{\"en-US\": \"Procurement\", \"zh-CN\": \"采购\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (308, 'global.finance', NULL, NULL, 'active', NULL, '{\"en-US\": \"Finance\", \"zh-CN\": \"财务\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (309, 'global.shipping', NULL, NULL, 'active', NULL, '{\"en-US\": \"Shipping\", \"zh-CN\": \"物流\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (310, 'global.packaging', NULL, NULL, 'active', NULL, '{\"en-US\": \"Packaging\", \"zh-CN\": \"包材\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (311, 'global.imports', NULL, NULL, 'active', NULL, '{\"en-US\": \"Imports\", \"zh-CN\": \"导入\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (312, 'global.systemSettings', NULL, NULL, 'active', NULL, '{\"en-US\": \"System Settings\", \"zh-CN\": \"系统设置\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (313, 'global.auditLogs', NULL, NULL, 'active', NULL, '{\"en-US\": \"Audit Logs\", \"zh-CN\": \"审计日志\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (314, 'global.jobs', NULL, NULL, 'active', NULL, '{\"en-US\": \"Jobs\", \"zh-CN\": \"任务\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (315, 'global.logs', NULL, NULL, 'active', NULL, '{\"en-US\": \"Logs\", \"zh-CN\": \"日志\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (316, 'global.dashboard', NULL, NULL, 'active', NULL, '{\"en-US\": \"Dashboard\", \"zh-CN\": \"仪表盘\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (317, 'global.settings', NULL, NULL, 'active', NULL, '{\"en-US\": \"Settings\", \"zh-CN\": \"设置\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (318, 'global.profile', NULL, NULL, 'active', NULL, '{\"en-US\": \"Profile\", \"zh-CN\": \"个人资料\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (319, 'global.language', NULL, NULL, 'active', NULL, '{\"en-US\": \"Language\", \"zh-CN\": \"语言\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (320, 'global.version', NULL, NULL, 'active', NULL, '{\"en-US\": \"Version\", \"zh-CN\": \"版本\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (321, 'global.print', NULL, NULL, 'active', NULL, '{\"en-US\": \"Print\", \"zh-CN\": \"打印\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (322, 'global.preview', NULL, NULL, 'active', NULL, '{\"en-US\": \"Preview\", \"zh-CN\": \"预览\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (323, 'global.apply', NULL, NULL, 'active', NULL, '{\"en-US\": \"Apply\", \"zh-CN\": \"应用\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (324, 'global.generate', NULL, NULL, 'active', NULL, '{\"en-US\": \"Generate\", \"zh-CN\": \"生成\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (325, 'global.sync', NULL, NULL, 'active', NULL, '{\"en-US\": \"Sync\", \"zh-CN\": \"同步\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (326, 'global.retrying', NULL, NULL, 'active', NULL, '{\"en-US\": \"Retrying\", \"zh-CN\": \"重试中\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (327, 'global.timeout', NULL, NULL, 'active', NULL, '{\"en-US\": \"Timeout\", \"zh-CN\": \"超时\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (328, 'global.networkError', NULL, NULL, 'active', NULL, '{\"en-US\": \"Network error\", \"zh-CN\": \"网络异常\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (329, 'global.serverError', NULL, NULL, 'active', NULL, '{\"en-US\": \"Server error\", \"zh-CN\": \"服务器异常\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (330, 'global.permissionDenied', NULL, NULL, 'active', NULL, '{\"en-US\": \"Permission denied\", \"zh-CN\": \"无权限\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (331, 'global.invalidFormat', NULL, NULL, 'active', NULL, '{\"en-US\": \"Invalid format\", \"zh-CN\": \"格式错误\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (332, 'global.requiredField', NULL, NULL, 'active', NULL, '{\"en-US\": \"Field required\", \"zh-CN\": \"字段必填\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (333, 'global.taxNo', NULL, NULL, 'active', NULL, '{\"en-US\": \"Tax ID\", \"zh-CN\": \"税号\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (334, 'global.bankAccount', NULL, NULL, 'active', NULL, '{\"en-US\": \"Bank Account\", \"zh-CN\": \"银行账户\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (335, 'global.bankName', NULL, NULL, 'active', NULL, '{\"en-US\": \"Bank Name\", \"zh-CN\": \"开户行\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (336, 'global.paymentTerms', NULL, NULL, 'active', NULL, '{\"en-US\": \"Payment Terms\", \"zh-CN\": \"付款条款\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (337, 'global.paymentStatus', NULL, NULL, 'active', NULL, '{\"en-US\": \"Payment Status\", \"zh-CN\": \"付款状态\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (338, 'global.paid', NULL, NULL, 'active', NULL, '{\"en-US\": \"Paid\", \"zh-CN\": \"已付款\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (339, 'global.unpaid', NULL, NULL, 'active', NULL, '{\"en-US\": \"Unpaid\", \"zh-CN\": \"未付款\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (340, 'global.partiallyPaid', NULL, NULL, 'active', NULL, '{\"en-US\": \"Partially Paid\", \"zh-CN\": \"部分付款\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (341, 'global.invoiceNo', NULL, NULL, 'active', NULL, '{\"en-US\": \"Invoice No.\", \"zh-CN\": \"发票号\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (342, 'global.dueDate', NULL, NULL, 'active', NULL, '{\"en-US\": \"Due Date\", \"zh-CN\": \"到期日\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (343, 'global.supplierCode', NULL, NULL, 'active', NULL, '{\"en-US\": \"Supplier Code\", \"zh-CN\": \"供应商编码\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (344, 'global.supplierType', NULL, NULL, 'active', NULL, '{\"en-US\": \"Supplier Type\", \"zh-CN\": \"供应商类型\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (345, 'global.customerCode', NULL, NULL, 'active', NULL, '{\"en-US\": \"Customer Code\", \"zh-CN\": \"客户编码\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (346, 'global.customerType', NULL, NULL, 'active', NULL, '{\"en-US\": \"Customer Type\", \"zh-CN\": \"客户类型\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (347, 'global.contactEmail', NULL, NULL, 'active', NULL, '{\"en-US\": \"Contact Email\", \"zh-CN\": \"联系邮箱\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (348, 'global.contactPhone', NULL, NULL, 'active', NULL, '{\"en-US\": \"Contact Phone\", \"zh-CN\": \"联系电话\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (349, 'global.shipmentNo', NULL, NULL, 'active', NULL, '{\"en-US\": \"Shipment No.\", \"zh-CN\": \"货件号\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (350, 'global.shipmentId', NULL, NULL, 'active', NULL, '{\"en-US\": \"Shipment ID\", \"zh-CN\": \"货件ID\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (351, 'global.shippingStatus', NULL, NULL, 'active', NULL, '{\"en-US\": \"Shipping Status\", \"zh-CN\": \"物流状态\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (352, 'global.picking', NULL, NULL, 'active', NULL, '{\"en-US\": \"Picking\", \"zh-CN\": \"拣货中\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (353, 'global.packed', NULL, NULL, 'active', NULL, '{\"en-US\": \"Packed\", \"zh-CN\": \"已打包\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (354, 'global.inTransit', NULL, NULL, 'active', NULL, '{\"en-US\": \"In Transit\", \"zh-CN\": \"运输中\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (355, 'global.delayed', NULL, NULL, 'active', NULL, '{\"en-US\": \"Delayed\", \"zh-CN\": \"已延迟\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (356, 'global.exception', NULL, NULL, 'active', NULL, '{\"en-US\": \"Exception\", \"zh-CN\": \"异常\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (357, 'global.cleared', NULL, NULL, 'active', NULL, '{\"en-US\": \"Cleared\", \"zh-CN\": \"已清关\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (358, 'global.receiving', NULL, NULL, 'active', NULL, '{\"en-US\": \"Receiving\", \"zh-CN\": \"收货中\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (359, 'global.receivedPartial', NULL, NULL, 'active', NULL, '{\"en-US\": \"Partially Received\", \"zh-CN\": \"部分收货\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (360, 'global.receivedFull', NULL, NULL, 'active', NULL, '{\"en-US\": \"Fully Received\", \"zh-CN\": \"全部收货\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (361, 'global.fbaShipmentId', NULL, NULL, 'active', NULL, '{\"en-US\": \"FBA Shipment ID\", \"zh-CN\": \"FBA货件号\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (362, 'global.fbaWarehouse', NULL, NULL, 'active', NULL, '{\"en-US\": \"FBA Warehouse\", \"zh-CN\": \"FBA仓\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (363, 'global.fc', NULL, NULL, 'active', NULL, '{\"en-US\": \"Fulfillment Center\", \"zh-CN\": \"运营中心\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (364, 'global.prepType', NULL, NULL, 'active', NULL, '{\"en-US\": \"Prep Type\", \"zh-CN\": \"预处理类型\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (365, 'global.labelType', NULL, NULL, 'active', NULL, '{\"en-US\": \"Label Type\", \"zh-CN\": \"标签类型\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (366, 'global.box', NULL, NULL, 'active', NULL, '{\"en-US\": \"Box\", \"zh-CN\": \"箱\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (367, 'global.case', NULL, NULL, 'active', NULL, '{\"en-US\": \"Case Pack\", \"zh-CN\": \"箱规\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (368, 'global.unitsPerCase', NULL, NULL, 'active', NULL, '{\"en-US\": \"Units per Case\", \"zh-CN\": \"每箱数量\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (369, 'global.unitsPerPallet', NULL, NULL, 'active', NULL, '{\"en-US\": \"Units per Pallet\", \"zh-CN\": \"每托数量\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (370, 'global.oversize', NULL, NULL, 'active', NULL, '{\"en-US\": \"Oversize\", \"zh-CN\": \"超大\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (371, 'global.standardSize', NULL, NULL, 'active', NULL, '{\"en-US\": \"Standard Size\", \"zh-CN\": \"标准尺寸\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (372, 'global.hazmat', NULL, NULL, 'active', NULL, '{\"en-US\": \"Hazmat\", \"zh-CN\": \"危险品\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (373, 'global.compliance', NULL, NULL, 'active', NULL, '{\"en-US\": \"Compliance\", \"zh-CN\": \"合规\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (374, 'global.inspection', NULL, NULL, 'active', NULL, '{\"en-US\": \"Inspection\", \"zh-CN\": \"质检\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (375, 'global.pass', NULL, NULL, 'active', NULL, '{\"en-US\": \"Pass\", \"zh-CN\": \"通过\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (376, 'global.fail', NULL, NULL, 'active', NULL, '{\"en-US\": \"Fail\", \"zh-CN\": \"不通过\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (377, 'global.procurementStatus', NULL, NULL, 'active', NULL, '{\"en-US\": \"Procurement Status\", \"zh-CN\": \"采购状态\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (378, 'global.poCreated', NULL, NULL, 'active', NULL, '{\"en-US\": \"Created\", \"zh-CN\": \"已创建\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (379, 'global.poApproved', NULL, NULL, 'active', NULL, '{\"en-US\": \"Approved\", \"zh-CN\": \"已审批\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (380, 'global.poOrdered', NULL, NULL, 'active', NULL, '{\"en-US\": \"Ordered\", \"zh-CN\": \"已下单\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (381, 'global.poPartialReceived', NULL, NULL, 'active', NULL, '{\"en-US\": \"Partially Received\", \"zh-CN\": \"部分到货\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (382, 'global.poReceived', NULL, NULL, 'active', NULL, '{\"en-US\": \"Received\", \"zh-CN\": \"已到货\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (383, 'global.poClosed', NULL, NULL, 'active', NULL, '{\"en-US\": \"Closed\", \"zh-CN\": \"已关闭\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (384, 'global.quote', NULL, NULL, 'active', NULL, '{\"en-US\": \"Quote\", \"zh-CN\": \"报价\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (385, 'global.quotationNo', NULL, NULL, 'active', NULL, '{\"en-US\": \"Quotation No.\", \"zh-CN\": \"报价单号\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (386, 'global.contract', NULL, NULL, 'active', NULL, '{\"en-US\": \"Contract\", \"zh-CN\": \"合同\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (387, 'global.contractNo', NULL, NULL, 'active', NULL, '{\"en-US\": \"Contract No.\", \"zh-CN\": \"合同号\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (388, 'global.settlementStatus', NULL, NULL, 'active', NULL, '{\"en-US\": \"Settlement Status\", \"zh-CN\": \"结算状态\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (389, 'global.settled', NULL, NULL, 'active', NULL, '{\"en-US\": \"Settled\", \"zh-CN\": \"已结算\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (390, 'global.unsettled', NULL, NULL, 'active', NULL, '{\"en-US\": \"Unsettled\", \"zh-CN\": \"未结算\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (391, 'global.partialSettled', NULL, NULL, 'active', NULL, '{\"en-US\": \"Partially Settled\", \"zh-CN\": \"部分结算\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (392, 'global.billing', NULL, NULL, 'active', NULL, '{\"en-US\": \"Billing\", \"zh-CN\": \"开票\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (393, 'global.billingStatus', NULL, NULL, 'active', NULL, '{\"en-US\": \"Billing Status\", \"zh-CN\": \"开票状态\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (394, 'global.billed', NULL, NULL, 'active', NULL, '{\"en-US\": \"Billed\", \"zh-CN\": \"已开票\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (395, 'global.unbilled', NULL, NULL, 'active', NULL, '{\"en-US\": \"Unbilled\", \"zh-CN\": \"未开票\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (396, 'global.partialBilled', NULL, NULL, 'active', NULL, '{\"en-US\": \"Partially Billed\", \"zh-CN\": \"部分开票\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (397, 'global.invoiceAmount', NULL, NULL, 'active', NULL, '{\"en-US\": \"Invoiced Amount\", \"zh-CN\": \"开票金额\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (398, 'global.taxAmount', NULL, NULL, 'active', NULL, '{\"en-US\": \"Tax Amount\", \"zh-CN\": \"税额\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (399, 'global.shippingFee', NULL, NULL, 'active', NULL, '{\"en-US\": \"Shipping Fee\", \"zh-CN\": \"运费\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (400, 'global.handlingFee', NULL, NULL, 'active', NULL, '{\"en-US\": \"Handling Fee\", \"zh-CN\": \"手续费\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (401, 'global.otherFee', NULL, NULL, 'active', NULL, '{\"en-US\": \"Other Fee\", \"zh-CN\": \"其他费用\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (402, 'global.costAllocation', NULL, NULL, 'active', NULL, '{\"en-US\": \"Cost Allocation\", \"zh-CN\": \"成本分摊\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (403, 'global.allocationRule', NULL, NULL, 'active', NULL, '{\"en-US\": \"Allocation Rule\", \"zh-CN\": \"分摊规则\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (404, 'global.account', NULL, NULL, 'active', NULL, '{\"en-US\": \"Account\", \"zh-CN\": \"账户\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (405, 'global.accountType', NULL, NULL, 'active', NULL, '{\"en-US\": \"Account Type\", \"zh-CN\": \"账户类型\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (406, 'global.accountNo', NULL, NULL, 'active', NULL, '{\"en-US\": \"Account No.\", \"zh-CN\": \"账号\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (407, 'global.accountBalance', NULL, NULL, 'active', NULL, '{\"en-US\": \"Account Balance\", \"zh-CN\": \"账户余额\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (408, 'global.debit', NULL, NULL, 'active', NULL, '{\"en-US\": \"Debit\", \"zh-CN\": \"借方\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (409, 'global.credit', NULL, NULL, 'active', NULL, '{\"en-US\": \"Credit\", \"zh-CN\": \"贷方\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (410, 'global.cash', NULL, NULL, 'active', NULL, '{\"en-US\": \"Cash\", \"zh-CN\": \"现金\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (411, 'global.bank', NULL, NULL, 'active', NULL, '{\"en-US\": \"Bank\", \"zh-CN\": \"银行\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (412, 'global.asset', NULL, NULL, 'active', NULL, '{\"en-US\": \"Asset\", \"zh-CN\": \"资产\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (413, 'global.liability', NULL, NULL, 'active', NULL, '{\"en-US\": \"Liability\", \"zh-CN\": \"负债\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (414, 'global.equity', NULL, NULL, 'active', NULL, '{\"en-US\": \"Equity\", \"zh-CN\": \"权益\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (415, 'global.revenue', NULL, NULL, 'active', NULL, '{\"en-US\": \"Revenue\", \"zh-CN\": \"收入\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (416, 'global.expense', NULL, NULL, 'active', NULL, '{\"en-US\": \"Expense\", \"zh-CN\": \"费用\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (417, 'global.cogs', NULL, NULL, 'active', NULL, '{\"en-US\": \"COGS\", \"zh-CN\": \"销售成本\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (418, 'global.grossProfit', NULL, NULL, 'active', NULL, '{\"en-US\": \"Gross Profit\", \"zh-CN\": \"毛利\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (419, 'global.netProfit', NULL, NULL, 'active', NULL, '{\"en-US\": \"Net Profit\", \"zh-CN\": \"净利\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (420, 'global.report', NULL, NULL, 'active', NULL, '{\"en-US\": \"Report\", \"zh-CN\": \"报表\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (421, 'global.reportPeriod', NULL, NULL, 'active', NULL, '{\"en-US\": \"Report Period\", \"zh-CN\": \"报表周期\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (422, 'global.monthly', NULL, NULL, 'active', NULL, '{\"en-US\": \"Monthly\", \"zh-CN\": \"月度\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (423, 'global.quarterly', NULL, NULL, 'active', NULL, '{\"en-US\": \"Quarterly\", \"zh-CN\": \"季度\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (424, 'global.yearly', NULL, NULL, 'active', NULL, '{\"en-US\": \"Yearly\", \"zh-CN\": \"年度\"}', '2026-01-19 15:38:20', '2026-01-19 15:38:20');
INSERT INTO `field_label` VALUES (491, 'product.images.title', 'product', 'images', 'active', NULL, '{\"en-US\": \"Image Manager\", \"zh-CN\": \"图片管理\"}', '2026-01-20 00:03:22', '2026-01-20 00:03:22');
INSERT INTO `field_label` VALUES (492, 'product.images.save', 'product', 'images', 'active', NULL, '{\"en-US\": \"Save\", \"zh-CN\": \"提交保存\"}', '2026-01-20 00:03:22', '2026-01-20 00:03:22');
INSERT INTO `field_label` VALUES (493, 'product.images.back', 'product', 'images', 'active', NULL, '{\"en-US\": \"Back to List\", \"zh-CN\": \"返回列表\"}', '2026-01-20 00:03:22', '2026-01-20 00:03:22');
INSERT INTO `field_label` VALUES (494, 'product.images.upload', 'product', 'images', 'active', NULL, '{\"en-US\": \"Upload Images\", \"zh-CN\": \"上传图片\"}', '2026-01-20 00:03:22', '2026-01-20 00:03:22');
INSERT INTO `field_label` VALUES (495, 'product.images.limit', 'product', 'images', 'active', NULL, '{\"en-US\": \"Up to 10 images, max 10MB each (jpg/jpeg/png/webp)\", \"zh-CN\": \"最多10张，单张不超过10MB（jpg/jpeg/png/webp）\"}', '2026-01-20 00:03:22', '2026-01-20 01:24:34');
INSERT INTO `field_label` VALUES (496, 'product.images.primary', 'product', 'images', 'active', NULL, '{\"en-US\": \"Primary\", \"zh-CN\": \"主图\"}', '2026-01-20 00:03:22', '2026-01-20 00:03:22');
INSERT INTO `field_label` VALUES (497, 'product.images.remove', 'product', 'images', 'active', NULL, '{\"en-US\": \"Remove\", \"zh-CN\": \"删除\"}', '2026-01-20 00:03:22', '2026-01-20 00:03:22');
INSERT INTO `field_label` VALUES (498, 'product.images.dirty', 'product', 'images', 'active', NULL, '{\"en-US\": \"Unsaved\", \"zh-CN\": \"未保存\"}', '2026-01-20 00:03:22', '2026-01-20 00:03:22');
INSERT INTO `field_label` VALUES (499, 'product.images.leaveConfirm', 'product', 'images', 'active', NULL, '{\"en-US\": \"You have unsaved changes. Leave page?\", \"zh-CN\": \"有未保存的图片更改，确认离开？\"}', '2026-01-20 00:03:22', '2026-01-20 00:03:22');
INSERT INTO `field_label` VALUES (500, 'product.images.warning', 'product', 'images', 'active', NULL, '{\"en-US\": \"Warning\", \"zh-CN\": \"提示\"}', '2026-01-20 00:03:22', '2026-01-20 00:03:22');
INSERT INTO `field_label` VALUES (501, 'product.images.batchUpload', 'product', 'images', 'active', NULL, '{\"en-US\": \"Batch Upload\", \"zh-CN\": \"批量上传\"}', '2026-01-20 00:16:46', '2026-01-20 00:16:46');
INSERT INTO `field_label` VALUES (502, 'product.images.singleUpload', 'product', 'images', 'active', NULL, '{\"en-US\": \"Upload\", \"zh-CN\": \"单张上传\"}', '2026-01-20 00:16:46', '2026-01-20 00:16:46');
INSERT INTO `field_label` VALUES (503, 'product.images.dragHint', 'product', 'images', 'active', NULL, '{\"en-US\": \"Drag to reorder, first image is primary\", \"zh-CN\": \"拖拽图片调整顺序，第一张为主图\"}', '2026-01-20 00:16:46', '2026-01-20 00:16:46');
INSERT INTO `field_label` VALUES (504, 'product.images.count', 'product', 'images', 'active', NULL, '{\"en-US\": \"Uploaded {count}/10\", \"zh-CN\": \"已上传 {count}/10\"}', '2026-01-20 00:16:46', '2026-01-20 00:16:46');
INSERT INTO `field_label` VALUES (529, 'product.images.limitSize', 'product', 'images', 'active', NULL, '{\"en-US\": \"Max 10MB per image (jpg/jpeg/png/webp)\", \"zh-CN\": \"单张不超过10MB（jpg/jpeg/png/webp）\"}', '2026-01-20 01:37:06', '2026-01-20 01:37:06');
INSERT INTO `field_label` VALUES (530, 'product.images.limitCount', 'product', 'images', 'active', NULL, '{\"en-US\": \"Up to 10 images\", \"zh-CN\": \"最多10张图片\"}', '2026-01-20 01:37:06', '2026-01-20 01:37:06');
INSERT INTO `field_label` VALUES (531, 'product.images.saveSuccess', 'product', 'images', 'active', NULL, '{\"en-US\": \"Saved\", \"zh-CN\": \"保存成功\"}', '2026-01-20 01:37:06', '2026-01-20 01:37:06');

-- ----------------------------
-- Table structure for inventory_balance
-- ----------------------------
DROP TABLE IF EXISTS `inventory_balance`;
CREATE TABLE `inventory_balance`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '库存余额ID',
  `sku_id` bigint UNSIGNED NOT NULL COMMENT 'SKU ID',
  `warehouse_id` bigint UNSIGNED NOT NULL COMMENT '仓库ID',
  `available_quantity` int UNSIGNED NOT NULL COMMENT '可用数量',
  `reserved_quantity` int UNSIGNED NOT NULL COMMENT '预留数量（已下单未发货）',
  `damaged_quantity` int UNSIGNED NOT NULL COMMENT '损坏数量',
  `purchasing_in_transit` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '采购在途库存',
  `pending_inspection` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '待检库存',
  `raw_material` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '原料库存',
  `pending_shipment` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '待出库存',
  `logistics_in_transit` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '物流在途库存',
  `sellable` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '可售库存',
  `returned` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '退货库存',
  `total_quantity` int UNSIGNED NOT NULL COMMENT '总数量',
  `last_movement_at` datetime NULL DEFAULT NULL COMMENT '最后流水时间',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_sku_warehouse`(`sku_id` ASC, `warehouse_id` ASC) USING BTREE,
  INDEX `idx_warehouse_id`(`warehouse_id` ASC) USING BTREE,
  INDEX `idx_sku_id`(`sku_id` ASC) USING BTREE,
  INDEX `idx_total_quantity`(`total_quantity` ASC) USING BTREE,
  INDEX `idx_available_quantity`(`available_quantity` ASC) USING BTREE,
  INDEX `idx_gmt_modified`(`gmt_modified` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '库存余额表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of inventory_balance
-- ----------------------------
INSERT INTO `inventory_balance` VALUES (1, 3, 1, 6, 0, 9, 150, 12, 10, 8, 7, 15, 0, 217, '2026-01-26 19:25:48', '2026-01-21 20:19:06', '2026-01-26 19:25:48');
INSERT INTO `inventory_balance` VALUES (2, 3, 2, 2, 0, 0, 10, 0, 10, 0, 0, 0, 0, 2, '2026-01-21 20:19:32', '2026-01-21 20:19:32', '2026-01-23 14:32:50');
INSERT INTO `inventory_balance` VALUES (3, 1, 1, 100, 0, 0, 10, 100, 6, 2, 2, 0, 0, 220, '2026-01-29 16:32:49', '2026-01-23 14:18:09', '2026-01-29 16:32:49');
INSERT INTO `inventory_balance` VALUES (4, 2, 1, 1, 0, 0, 10, 1, 10, 0, 0, 0, 0, 22, '2026-01-24 12:27:13', '2026-01-23 14:18:09', '2026-01-24 12:27:13');

-- ----------------------------
-- Table structure for inventory_movement
-- ----------------------------
DROP TABLE IF EXISTS `inventory_movement`;
CREATE TABLE `inventory_movement`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '流水ID',
  `trace_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '追踪ID（关联审计日志）',
  `sku_id` bigint UNSIGNED NOT NULL COMMENT 'SKU ID',
  `warehouse_id` bigint UNSIGNED NOT NULL COMMENT '仓库ID',
  `movement_type` enum('PURCHASE_RECEIPT','SALES_SHIPMENT','STOCK_TAKE_ADJUSTMENT','MANUAL_ADJUSTMENT','DAMAGE_WRITE_OFF','RETURN_RECEIPT','TRANSFER_OUT','TRANSFER_IN','PURCHASE_SHIP','WAREHOUSE_RECEIVE','INSPECTION_PASS','INSPECTION_FAIL','ASSEMBLY_COMPLETE','LOGISTICS_SHIP','PLATFORM_RECEIVE','RETURN_INSPECT','SHIPMENT_LOCK','SHIPMENT_PACK','SHIPMENT_SHIP','SHIPMENT_ROLLBACK') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '流水类型',
  `reference_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '参考单据类型（如PurchaseOrder, SalesOrder）',
  `reference_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '参考单据ID',
  `reference_number` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '参考单据号',
  `quantity` int NOT NULL COMMENT '变动数量（正数=入库，负数=出库）',
  `before_available` int UNSIGNED NOT NULL COMMENT '变动前可用数量',
  `after_available` int UNSIGNED NOT NULL COMMENT '变动后可用数量',
  `before_reserved` int UNSIGNED NOT NULL COMMENT '变动前预留数量',
  `after_reserved` int UNSIGNED NOT NULL COMMENT '变动后预留数量',
  `before_damaged` int UNSIGNED NOT NULL COMMENT '变动前损坏数量',
  `after_damaged` int UNSIGNED NOT NULL COMMENT '变动后损坏数量',
  `unit_cost` decimal(12, 4) NULL DEFAULT NULL COMMENT '单位成本',
  `total_cost` decimal(12, 4) NULL DEFAULT NULL COMMENT '总成本',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '备注',
  `operator_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '操作人ID',
  `operated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_trace_id`(`trace_id` ASC) USING BTREE,
  INDEX `idx_sku_id`(`sku_id` ASC) USING BTREE,
  INDEX `idx_warehouse_id`(`warehouse_id` ASC) USING BTREE,
  INDEX `idx_movement_type`(`movement_type` ASC) USING BTREE,
  INDEX `idx_reference_type_reference_id`(`reference_type` ASC, `reference_id` ASC) USING BTREE,
  INDEX `idx_operated_at`(`operated_at` ASC) USING BTREE,
  INDEX `idx_operator_id`(`operator_id` ASC) USING BTREE,
  INDEX `idx_sku_warehouse_operated`(`sku_id` ASC, `warehouse_id` ASC, `operated_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 34 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '库存流水表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of inventory_movement
-- ----------------------------
INSERT INTO `inventory_movement` VALUES (1, '', 3, 1, 'PURCHASE_RECEIPT', 'PURCHASE_ORDER', 40, 'PO202601212000498500', 1, 0, 1, 0, 0, 0, 0, 1.0000, 1.0000, '', 8, '2026-01-21 20:04:13', '2026-01-21 20:04:13', '2026-01-21 20:04:13');
INSERT INTO `inventory_movement` VALUES (2, '6c40ca4a-1e0a-474f-babe-48605605466d', 3, 1, 'PURCHASE_RECEIPT', NULL, NULL, NULL, 10, 0, 10, 0, 0, 0, 0, 15.5000, 155.0000, '���Բɹ�����', 8, '2026-01-21 20:19:06', '2026-01-21 20:19:06', '2026-01-21 20:19:06');
INSERT INTO `inventory_movement` VALUES (3, '839a00d1-5fb9-46e3-a47e-74fde72fbbe5', 3, 1, 'SALES_SHIPMENT', NULL, NULL, NULL, -3, 10, 7, 0, 0, 0, 0, NULL, NULL, '�������۳���', 8, '2026-01-21 20:19:25', '2026-01-21 20:19:25', '2026-01-21 20:19:25');
INSERT INTO `inventory_movement` VALUES (4, 'df5cb7bc-db75-40a9-b4cd-e47c899a70c4', 3, 1, 'TRANSFER_OUT', NULL, NULL, NULL, -2, 7, 5, 0, 0, 0, 0, 15.5000, 31.0000, '���Ե���', 8, '2026-01-21 20:19:32', '2026-01-21 20:19:32', '2026-01-21 20:19:32');
INSERT INTO `inventory_movement` VALUES (5, '646d6adb-6630-432d-ab06-626af217625e', 3, 2, 'TRANSFER_IN', NULL, NULL, NULL, 2, 0, 2, 0, 0, 0, 0, 15.5000, 31.0000, '���Ե���', 8, '2026-01-21 20:19:32', '2026-01-21 20:19:32', '2026-01-21 20:19:32');
INSERT INTO `inventory_movement` VALUES (6, '4b1f69c5-fdc8-474e-b2fd-4198027f97c6', 3, 1, 'PURCHASE_SHIP', NULL, NULL, NULL, 100, 5, 5, 0, 0, 0, 0, NULL, NULL, '��Ӧ�̷�������', NULL, '2026-01-22 12:15:20', '2026-01-22 12:15:20', '2026-01-22 12:15:20');
INSERT INTO `inventory_movement` VALUES (7, 'ad329923-75d3-4e74-9125-959540bab256', 3, 1, 'PURCHASE_SHIP', NULL, NULL, NULL, 100, 5, 5, 0, 0, 0, 0, NULL, NULL, '��Ӧ�̷�������', NULL, '2026-01-22 12:16:06', '2026-01-22 12:16:06', '2026-01-22 12:16:06');
INSERT INTO `inventory_movement` VALUES (8, 'db664dbf-150b-40f7-a61c-79098822db5f', 3, 1, 'WAREHOUSE_RECEIVE', NULL, NULL, NULL, 50, 5, 5, 0, 0, 0, 0, NULL, NULL, '�����ջ�����', NULL, '2026-01-22 12:16:36', '2026-01-22 12:16:36', '2026-01-22 12:16:36');
INSERT INTO `inventory_movement` VALUES (9, 'e9920077-883e-4735-9d2a-3b82cad028b6', 3, 1, 'INSPECTION_PASS', NULL, NULL, NULL, 40, 5, 5, 0, 0, 0, 0, NULL, NULL, '�ʼ�ͨ������', NULL, '2026-01-22 12:16:36', '2026-01-22 12:16:36', '2026-01-22 12:16:36');
INSERT INTO `inventory_movement` VALUES (10, 'bed3de1e-336a-459d-8d83-087cc0bb76b4', 3, 1, 'INSPECTION_FAIL', NULL, NULL, NULL, 5, 5, 5, 0, 0, 0, 5, NULL, NULL, '�ʼ첻�ϸ�', NULL, '2026-01-22 12:16:36', '2026-01-22 12:16:36', '2026-01-22 12:16:36');
INSERT INTO `inventory_movement` VALUES (11, 'b86fe657-6771-45bc-8b31-d6cd994596c0', 3, 1, 'ASSEMBLY_COMPLETE', NULL, NULL, NULL, 30, 5, 5, 0, 0, 5, 5, NULL, NULL, '��װ����', NULL, '2026-01-22 12:17:15', '2026-01-22 12:17:15', '2026-01-22 12:17:15');
INSERT INTO `inventory_movement` VALUES (12, '1bf55d47-d194-4e7d-bd75-899b7808c60c', 3, 1, 'LOGISTICS_SHIP', NULL, NULL, NULL, 20, 5, 5, 0, 0, 5, 5, NULL, NULL, '����FBA', NULL, '2026-01-22 12:17:15', '2026-01-22 12:17:15', '2026-01-22 12:17:15');
INSERT INTO `inventory_movement` VALUES (13, '3bf82edf-e67a-4538-92a8-54154f94e332', 3, 1, 'PLATFORM_RECEIVE', NULL, NULL, NULL, 15, 5, 5, 0, 0, 5, 5, NULL, NULL, 'FBA�ϼ�', NULL, '2026-01-22 12:17:15', '2026-01-22 12:17:15', '2026-01-22 12:17:15');
INSERT INTO `inventory_movement` VALUES (14, '1d0bd9af-b6ab-44c5-881d-8c3954382241', 3, 1, 'RETURN_RECEIPT', NULL, NULL, NULL, 10, 5, 5, 0, 0, 5, 5, NULL, NULL, '�ͻ��˻�', NULL, '2026-01-22 12:17:35', '2026-01-22 12:17:35', '2026-01-22 12:17:35');
INSERT INTO `inventory_movement` VALUES (15, '327bd15b-9584-43d0-bc72-907aca19a3a4', 3, 1, 'RETURN_INSPECT', NULL, NULL, NULL, 6, 5, 5, 0, 0, 5, 5, NULL, NULL, '�˻��ʼ�', NULL, '2026-01-22 12:17:35', '2026-01-22 12:17:35', '2026-01-22 12:17:35');
INSERT INTO `inventory_movement` VALUES (16, '327bd15b-9584-43d0-bc72-907aca19a3a4', 3, 1, 'RETURN_INSPECT', NULL, NULL, NULL, -4, 5, 5, 0, 0, 5, 9, NULL, NULL, '�˻��ʼ� - 质检不合格', NULL, '2026-01-22 12:17:35', '2026-01-22 12:17:35', '2026-01-22 12:17:35');
INSERT INTO `inventory_movement` VALUES (17, '', 1, 1, 'PURCHASE_RECEIPT', 'PURCHASE_ORDER', 44, 'PO202601221716443200', 100, 0, 100, 0, 0, 0, 0, 100.0000, 10000.0000, '', 8, '2026-01-22 17:17:29', '2026-01-22 17:17:29', '2026-01-22 17:17:29');
INSERT INTO `inventory_movement` VALUES (18, '', 2, 1, 'PURCHASE_RECEIPT', 'PURCHASE_ORDER', 44, 'PO202601221716443200', 1, 0, 1, 0, 0, 0, 0, 1.0000, 1.0000, '', 8, '2026-01-22 17:17:29', '2026-01-22 17:17:29', '2026-01-22 17:17:29');
INSERT INTO `inventory_movement` VALUES (19, 'dc2713a5-e9f3-4dac-83fb-31b9bed821ab', 3, 1, 'PURCHASE_RECEIPT', 'PURCHASE_ORDER', 80, 'PO20260123141628c78db3ce', 1, 5, 6, 0, 0, 9, 9, 1.0000, 1.0000, '采购入库: PO20260123141628c78db3ce', 8, '2026-01-23 14:17:44', '2026-01-23 14:17:44', '2026-01-23 14:17:44');
INSERT INTO `inventory_movement` VALUES (20, 'd3544f7e-5399-4604-b7d7-8d11932fccc5', 1, 1, 'PURCHASE_RECEIPT', 'PURCHASE_ORDER', 81, 'PO202601231416282515eabc', 100, 0, 100, 0, 0, 0, 0, 100.0000, 10000.0000, '采购入库: PO202601231416282515eabc', 8, '2026-01-23 14:18:09', '2026-01-23 14:18:09', '2026-01-23 14:18:09');
INSERT INTO `inventory_movement` VALUES (21, '4080451e-149c-467b-b8c3-5314502a2ee8', 2, 1, 'PURCHASE_RECEIPT', 'PURCHASE_ORDER', 81, 'PO202601231416282515eabc', 1, 0, 1, 0, 0, 0, 0, 1.0000, 1.0000, '采购入库: PO202601231416282515eabc', 8, '2026-01-23 14:18:09', '2026-01-23 14:18:09', '2026-01-23 14:18:09');
INSERT INTO `inventory_movement` VALUES (22, '16b93504-2660-414e-93e5-cb206cfc6092', 1, 1, 'ASSEMBLY_COMPLETE', NULL, NULL, NULL, 1, 100, 100, 0, 0, 0, 0, NULL, NULL, NULL, 8, '2026-01-23 14:33:12', '2026-01-23 14:33:12', '2026-01-23 14:33:12');
INSERT INTO `inventory_movement` VALUES (23, '446c9990-f532-4ab9-94c0-156ea61a73ca', 3, 1, 'PURCHASE_SHIP', 'PURCHASE_ORDER', 79, 'PO202601231416094da24355', 1, 6, 6, 0, 0, 9, 9, 1.0000, 1.0000, '采购发货入在途: PO202601231416094da24355', 8, '2026-01-24 12:26:56', '2026-01-24 12:26:56', '2026-01-24 12:26:56');
INSERT INTO `inventory_movement` VALUES (24, '8bfc28b5-a32c-4e73-b228-68d40649763b', 3, 1, 'WAREHOUSE_RECEIVE', 'PURCHASE_ORDER', 79, 'PO202601231416094da24355', 1, 6, 6, 0, 0, 9, 9, 1.0000, 1.0000, '仓库收货入待检: PO202601231416094da24355', 8, '2026-01-24 12:27:01', '2026-01-24 12:27:01', '2026-01-24 12:27:01');
INSERT INTO `inventory_movement` VALUES (25, '6a3855f6-dc87-4439-bb10-1a654b59a87c', 1, 1, 'PURCHASE_SHIP', 'PURCHASE_ORDER', 78, 'PO20260123141609864d37d4', 100, 100, 100, 0, 0, 0, 0, 100.0000, 10000.0000, '采购发货入在途: PO20260123141609864d37d4', 8, '2026-01-24 12:27:08', '2026-01-24 12:27:08', '2026-01-24 12:27:08');
INSERT INTO `inventory_movement` VALUES (26, 'ef690fa8-1c5d-4d75-80c5-e1d195d542ec', 2, 1, 'PURCHASE_SHIP', 'PURCHASE_ORDER', 78, 'PO20260123141609864d37d4', 1, 1, 1, 0, 0, 0, 0, 1.0000, 1.0000, '采购发货入在途: PO20260123141609864d37d4', 8, '2026-01-24 12:27:08', '2026-01-24 12:27:08', '2026-01-24 12:27:08');
INSERT INTO `inventory_movement` VALUES (27, '378f371c-efc3-4c9e-b003-07ce5278a9e7', 1, 1, 'WAREHOUSE_RECEIVE', 'PURCHASE_ORDER', 78, 'PO20260123141609864d37d4', 100, 100, 100, 0, 0, 0, 0, 100.0000, 10000.0000, '仓库收货入待检: PO20260123141609864d37d4', 8, '2026-01-24 12:27:13', '2026-01-24 12:27:13', '2026-01-24 12:27:13');
INSERT INTO `inventory_movement` VALUES (28, 'c7ec20a6-2865-457a-83e0-d750bfa5f635', 2, 1, 'WAREHOUSE_RECEIVE', 'PURCHASE_ORDER', 78, 'PO20260123141609864d37d4', 1, 1, 1, 0, 0, 0, 0, 1.0000, 1.0000, '仓库收货入待检: PO20260123141609864d37d4', 8, '2026-01-24 12:27:13', '2026-01-24 12:27:13', '2026-01-24 12:27:13');
INSERT INTO `inventory_movement` VALUES (29, 'b43bbecc-da79-4ebf-8a32-383508f9a057', 1, 1, 'ASSEMBLY_COMPLETE', NULL, NULL, NULL, 1, 100, 100, 0, 0, 0, 0, NULL, NULL, NULL, 8, '2026-01-24 12:27:47', '2026-01-24 12:27:47', '2026-01-24 12:27:47');
INSERT INTO `inventory_movement` VALUES (30, '0e062f45-0734-4b72-9834-b2ea41ad6f26', 3, 1, 'SHIPMENT_SHIP', 'SHIPMENT', 5, 'SH2026012619245407779', 2, 6, 6, 0, 0, 9, 9, 0.0000, 0.0000, NULL, 8, '2026-01-26 19:25:48', '2026-01-26 19:25:48', '2026-01-26 19:25:48');
INSERT INTO `inventory_movement` VALUES (31, 'ef83fe97-86ec-4f26-a138-7cf4fce209f8', 1, 1, 'SHIPMENT_SHIP', 'SHIPMENT', 8, 'SH2026012814545113557', 2, 100, 100, 0, 0, 0, 0, 0.0000, 0.0000, NULL, 8, '2026-01-28 14:55:15', '2026-01-28 14:55:15', '2026-01-28 14:55:15');
INSERT INTO `inventory_movement` VALUES (32, 'bfdf2f43-6b28-4fb3-bf95-b4003706f12c', 1, 1, 'ASSEMBLY_COMPLETE', NULL, NULL, NULL, 1, 100, 100, 0, 0, 0, 0, NULL, NULL, NULL, 8, '2026-01-28 14:57:20', '2026-01-28 14:57:20', '2026-01-28 14:57:20');
INSERT INTO `inventory_movement` VALUES (33, '41a450de-f96f-423b-b443-6f985825edcd', 1, 1, 'ASSEMBLY_COMPLETE', NULL, NULL, NULL, 1, 100, 100, 0, 0, 0, 0, NULL, NULL, NULL, 8, '2026-01-29 16:32:49', '2026-01-29 16:32:49', '2026-01-29 16:32:49');

-- ----------------------------
-- Table structure for job_run
-- ----------------------------
DROP TABLE IF EXISTS `job_run`;
CREATE TABLE `job_run`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '任务执行ID',
  `trace_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '追踪ID（全链路追踪）',
  `job_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '任务类型（IMPORT/RECONCILE/REPLENISH_CALC/COST_CALC等）',
  `job_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '任务名称（描述性名称）',
  `status` enum('PENDING','RUNNING','SUCCESS','FAILED','CANCELLED') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'PENDING' COMMENT '执行状态',
  `started_at` datetime NULL DEFAULT NULL COMMENT '开始时间',
  `finished_at` datetime NULL DEFAULT NULL COMMENT '结束时间',
  `duration_ms` int UNSIGNED NULL DEFAULT NULL COMMENT '耗时（毫秒）',
  `total_rows` int UNSIGNED NULL DEFAULT NULL COMMENT '总行数',
  `success_rows` int UNSIGNED NULL DEFAULT NULL COMMENT '成功行数',
  `failed_rows` int UNSIGNED NULL DEFAULT NULL COMMENT '失败行数',
  `input_summary` json NULL COMMENT '输入摘要（文件名、参数、配置等）',
  `output_summary` json NULL COMMENT '输出摘要（结果统计、关键指标）',
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '错误信息',
  `created_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '创建人ID（手动触发时记录）',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_trace_id`(`trace_id` ASC) USING BTREE,
  INDEX `idx_job_type`(`job_type` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_created_by`(`created_by` ASC) USING BTREE,
  INDEX `idx_gmt_create`(`gmt_create` ASC) USING BTREE,
  INDEX `idx_status_job_type`(`status` ASC, `job_type` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '任务执行记录表（导入、对账、计算等长任务）' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of job_run
-- ----------------------------

-- ----------------------------
-- Table structure for job_run_item
-- ----------------------------
DROP TABLE IF EXISTS `job_run_item`;
CREATE TABLE `job_run_item`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '明细ID',
  `job_run_id` bigint UNSIGNED NOT NULL COMMENT '任务执行ID',
  `row_number` int UNSIGNED NULL DEFAULT NULL COMMENT '行号（CSV行号、数据行号等）',
  `status` enum('SUCCESS','FAILED','SKIPPED') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '处理状态',
  `item_data` json NULL COMMENT '行数据摘要（关键字段，不存完整数据）',
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '错误信息',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_job_run_id`(`job_run_id` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_job_run_id_status`(`job_run_id` ASC, `status` ASC) USING BTREE,
  INDEX `idx_gmt_create`(`gmt_create` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '任务执行明细表（失败行、错误行、跳过行记录）' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of job_run_item
-- ----------------------------

-- ----------------------------
-- Table structure for logistics_provider
-- ----------------------------
DROP TABLE IF EXISTS `logistics_provider`;
CREATE TABLE `logistics_provider`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `provider_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '供应商代码',
  `provider_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '供应商名称',
  `provider_type` enum('FREIGHT_FORWARDER','COURIER','SHIPPING_LINE','AIRLINE') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '供应商类型: 货代/快递/船公司/航空',
  `service_types` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '服务类型(EXPRESS,AIR,SEA,RAIL),逗号分隔',
  `contact_person` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '联系人',
  `contact_phone` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '联系电话',
  `contact_email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '联系邮箱',
  `website` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '网站',
  `country` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '国家',
  `city` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '城市',
  `address` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '地址',
  `account_number` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '客户账号',
  `credit_days` int NULL DEFAULT 0 COMMENT '账期天数',
  `status` enum('ACTIVE','INACTIVE') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '备注',
  `created_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '更新人',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_provider_code`(`provider_code` ASC) USING BTREE,
  INDEX `idx_provider_type`(`provider_type` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 14 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '物流供应商表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of logistics_provider
-- ----------------------------
INSERT INTO `logistics_provider` VALUES (1, 'ju_fang', '九方物流', 'FREIGHT_FORWARDER', '', '', '', '', NULL, NULL, NULL, '', '', 0, 'ACTIVE', '', 8, 8, '2026-01-27 16:18:11', '2026-01-27 16:18:11');
INSERT INTO `logistics_provider` VALUES (2, 'DHL-CN', 'DHL中国', 'COURIER', 'EXPRESS,AIR', '张三', '400-810-8000', 'service@dhl.com.cn', NULL, '中国', '上海', NULL, NULL, 0, 'ACTIVE', NULL, NULL, NULL, '2026-01-28 14:07:36', '2026-01-28 14:07:36');
INSERT INTO `logistics_provider` VALUES (3, 'FEDEX-CN', 'FedEx中国', 'COURIER', 'EXPRESS,AIR', '李四', '400-886-1888', 'service@fedex.com.cn', NULL, '中国', '上海', NULL, NULL, 0, 'ACTIVE', NULL, NULL, NULL, '2026-01-28 14:07:36', '2026-01-28 14:07:36');
INSERT INTO `logistics_provider` VALUES (4, 'UPS-CN', 'UPS中国', 'COURIER', 'EXPRESS,AIR', '王五', '400-820-8388', 'service@ups.com.cn', NULL, '中国', '上海', NULL, NULL, 0, 'ACTIVE', NULL, NULL, NULL, '2026-01-28 14:07:36', '2026-01-28 14:07:36');
INSERT INTO `logistics_provider` VALUES (5, 'SF-EXPRESS', '顺丰速运', 'COURIER', 'EXPRESS,AIR', '赵六', '95338', 'service@sf-express.com', NULL, '中国', '深圳', NULL, NULL, 0, 'ACTIVE', NULL, NULL, NULL, '2026-01-28 14:07:36', '2026-01-28 14:07:36');
INSERT INTO `logistics_provider` VALUES (6, 'FREIGHT-FWD-01', '环球货运代理', 'FREIGHT_FORWARDER', 'AIR,SEA', '钱七', '021-12345678', 'info@global-freight.com', NULL, '中国', '上海', NULL, NULL, 0, 'ACTIVE', NULL, NULL, NULL, '2026-01-28 14:07:36', '2026-01-28 14:07:36');
INSERT INTO `logistics_provider` VALUES (7, 'FREIGHT-FWD-02', '中远海运物流', 'FREIGHT_FORWARDER', 'SEA,RAIL', '孙八', '0755-88888888', 'service@cosco-logistics.com', NULL, '中国', '深圳', NULL, NULL, 0, 'ACTIVE', NULL, NULL, NULL, '2026-01-28 14:07:36', '2026-01-28 14:07:36');
INSERT INTO `logistics_provider` VALUES (8, 'FREIGHT-FWD-03', '嘉里物流', 'FREIGHT_FORWARDER', 'AIR,SEA,RAIL,TRUCK', '周九', '400-820-3333', 'service@kerrylogistics.com', NULL, '中国', '上海', NULL, NULL, 0, 'ACTIVE', NULL, NULL, NULL, '2026-01-28 14:07:36', '2026-01-28 14:07:36');
INSERT INTO `logistics_provider` VALUES (9, 'MAERSK', '马士基航运', 'SHIPPING_LINE', 'SEA', '吴十', '400-120-6888', 'service@maersk.com.cn', NULL, '中国', '上海', NULL, NULL, 0, 'ACTIVE', NULL, NULL, NULL, '2026-01-28 14:07:36', '2026-01-28 14:07:36');
INSERT INTO `logistics_provider` VALUES (10, 'MSC', '地中海航运', 'SHIPPING_LINE', 'SEA', '郑十一', '021-63303000', 'service@msc.com.cn', NULL, '中国', '上海', NULL, NULL, 0, 'ACTIVE', NULL, NULL, NULL, '2026-01-28 14:07:36', '2026-01-28 14:07:36');
INSERT INTO `logistics_provider` VALUES (11, 'COSCO', '中远海运集运', 'SHIPPING_LINE', 'SEA', '林十二', '400-820-8888', 'service@coscoshipping.com', NULL, '中国', '上海', NULL, NULL, 0, 'ACTIVE', NULL, NULL, NULL, '2026-01-28 14:07:36', '2026-01-28 14:07:36');
INSERT INTO `logistics_provider` VALUES (12, 'CA-CARGO', '国航货运', 'AIRLINE', 'AIR', '陈十三', '010-95583', 'cargo@airchina.com', NULL, '中国', '北京', NULL, NULL, 0, 'ACTIVE', NULL, NULL, NULL, '2026-01-28 14:07:36', '2026-01-28 14:07:36');
INSERT INTO `logistics_provider` VALUES (13, 'CZ-CARGO', '南航货运', 'AIRLINE', 'AIR', '黄十四', '020-95539', 'cargo@csair.com', NULL, '中国', '广州', NULL, NULL, 0, 'ACTIVE', NULL, NULL, NULL, '2026-01-28 14:07:36', '2026-01-28 14:07:36');

-- ----------------------------
-- Table structure for logistics_service
-- ----------------------------
DROP TABLE IF EXISTS `logistics_service`;
CREATE TABLE `logistics_service`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `service_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '服务代码（如：SLOW_SHIP、FAST_SHIP）',
  `service_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '服务名称（如：慢船、快船、美森快船）',
  `transport_mode` enum('EXPRESS','AIR','SEA','RAIL','TRUCK') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '运输方式',
  `destination_region` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '目的地站点/国家（如：美国、欧洲、日本）',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '服务描述',
  `status` enum('ACTIVE','INACTIVE') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
  `created_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '创建人ID',
  `updated_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '更新人ID',
  `gmt_create` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `service_code`(`service_code` ASC) USING BTREE,
  INDEX `idx_transport_mode`(`transport_mode` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 27 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '物流服务表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of logistics_service
-- ----------------------------
INSERT INTO `logistics_service` VALUES (1, 'SEA_SLOW_US', '美国慢船', 'SEA', '美国', '经济型海运服务，时效较慢但价格优惠', 'ACTIVE', NULL, NULL, '2026-01-27 19:03:58', '2026-01-27 19:03:58');
INSERT INTO `logistics_service` VALUES (2, 'SEA_FAST_US', '美国快船', 'SEA', '美国', '快速海运服务，时效较快', 'ACTIVE', NULL, NULL, '2026-01-27 19:03:58', '2026-01-27 19:03:58');
INSERT INTO `logistics_service` VALUES (3, 'SEA_MATSON_US', '美森快船', 'SEA', '美国', 'Matson快船服务，时效最快', 'ACTIVE', NULL, NULL, '2026-01-27 19:03:58', '2026-01-27 19:03:58');
INSERT INTO `logistics_service` VALUES (4, 'SEA_SLOW_EU', '欧洲慢船', 'SEA', '欧洲', '经济型海运服务到欧洲', 'ACTIVE', NULL, NULL, '2026-01-27 19:03:58', '2026-01-27 19:03:58');
INSERT INTO `logistics_service` VALUES (5, 'AIR_STANDARD', '标准空运', 'AIR', '全球', '标准空运服务', 'ACTIVE', NULL, NULL, '2026-01-27 19:03:58', '2026-01-27 19:03:58');
INSERT INTO `logistics_service` VALUES (6, 'AIR_EXPRESS', '快速空运', 'AIR', '全球', '加急空运服务', 'ACTIVE', NULL, NULL, '2026-01-27 19:03:58', '2026-01-27 19:03:58');
INSERT INTO `logistics_service` VALUES (7, 'EXPRESS_DHL', 'DHL快递', 'EXPRESS', '全球', 'DHL国际快递服务', 'ACTIVE', NULL, NULL, '2026-01-27 19:03:58', '2026-01-27 19:03:58');
INSERT INTO `logistics_service` VALUES (8, 'EXPRESS_FEDEX', 'FedEx快递', 'EXPRESS', '全球', 'FedEx国际快递服务', 'ACTIVE', NULL, NULL, '2026-01-27 19:03:58', '2026-01-27 19:03:58');
INSERT INTO `logistics_service` VALUES (9, 'EXPRESS-DOMESTIC', '国内快递', 'EXPRESS', '中国大陆', '国内标准快递服务，1-3天送达', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');
INSERT INTO `logistics_service` VALUES (10, 'EXPRESS-INTL-STD', '国际标准快递', 'EXPRESS', '全球', '国际标准快递，3-7个工作日', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');
INSERT INTO `logistics_service` VALUES (11, 'EXPRESS-INTL-PRIORITY', '国际优先快递', 'EXPRESS', '全球', '国际优先快递，2-4个工作日', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');
INSERT INTO `logistics_service` VALUES (12, 'EXPRESS-USA', '美国快递专线', 'EXPRESS', '美国', '美国快递专线，3-5个工作日', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');
INSERT INTO `logistics_service` VALUES (13, 'EXPRESS-EU', '欧洲快递专线', 'EXPRESS', '欧洲', '欧洲快递专线，4-7个工作日', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');
INSERT INTO `logistics_service` VALUES (14, 'AIR-STANDARD', '标准空运', 'AIR', '全球', '标准空运服务，5-10个工作日', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');
INSERT INTO `logistics_service` VALUES (15, 'AIR-FAST', '快速空运', 'AIR', '全球', '快速空运服务，3-5个工作日', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');
INSERT INTO `logistics_service` VALUES (16, 'AIR-USA', '美国空运专线', 'AIR', '美国', '美国空运专线，4-6个工作日', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');
INSERT INTO `logistics_service` VALUES (17, 'AIR-EU', '欧洲空运专线', 'AIR', '欧洲', '欧洲空运专线，5-8个工作日', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');
INSERT INTO `logistics_service` VALUES (18, 'SEA-FCL', '整柜海运', 'SEA', '全球', '整柜海运服务，20-40天', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');
INSERT INTO `logistics_service` VALUES (19, 'SEA-LCL', '拼箱海运', 'SEA', '全球', '拼箱海运服务，25-45天', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');
INSERT INTO `logistics_service` VALUES (20, 'SEA-USA-WEST', '美西海运', 'SEA', '美国西海岸', '美国西海岸海运，15-20天', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');
INSERT INTO `logistics_service` VALUES (21, 'SEA-USA-EAST', '美东海运', 'SEA', '美国东海岸', '美国东海岸海运，25-30天', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');
INSERT INTO `logistics_service` VALUES (22, 'SEA-EU', '欧洲海运', 'SEA', '欧洲', '欧洲海运服务，30-40天', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');
INSERT INTO `logistics_service` VALUES (23, 'RAIL-EU', '中欧班列', 'RAIL', '欧洲', '中欧铁路班列，15-20天', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');
INSERT INTO `logistics_service` VALUES (24, 'RAIL-RUSSIA', '中俄班列', 'RAIL', '俄罗斯', '中俄铁路班列，10-15天', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');
INSERT INTO `logistics_service` VALUES (25, 'TRUCK-DOMESTIC', '国内卡车运输', 'TRUCK', '中国大陆', '国内卡车整车/零担运输', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');
INSERT INTO `logistics_service` VALUES (26, 'TRUCK-CROSS-BORDER', '跨境卡车运输', 'TRUCK', '东南亚', '东南亚跨境卡车运输', 'ACTIVE', NULL, NULL, '2026-01-28 14:07:44', '2026-01-28 14:07:44');

-- ----------------------------
-- Table structure for menu
-- ----------------------------
DROP TABLE IF EXISTS `menu`;
CREATE TABLE `menu`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '菜单名称',
  `title_en` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '菜单英文名',
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '菜单编码',
  `parent_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '父级菜单ID',
  `path` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '前端路由',
  `icon` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '图标',
  `component` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '组件',
  `sort` int UNSIGNED NOT NULL COMMENT '排序',
  `is_hidden` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否隐藏',
  `permission_code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '权限编码',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_code`(`code` ASC) USING BTREE,
  INDEX `idx_parent_id`(`parent_id` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  CONSTRAINT `chk_menu_status` CHECK (`status` in (_gbk'ACTIVE',_gbk'DISABLED'))
) ENGINE = InnoDB AUTO_INCREMENT = 36 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统菜单' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of menu
-- ----------------------------
INSERT INTO `menu` VALUES (1, '系统管理', 'System', 'SYSTEM', NULL, NULL, 'Setting', NULL, 100, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 20:17:48');
INSERT INTO `menu` VALUES (2, '商品管理', 'Product', 'PRODUCT', NULL, NULL, 'Goods', NULL, 1, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 23:32:55');
INSERT INTO `menu` VALUES (3, '库存管理', 'Inventory', 'INVENTORY', NULL, NULL, 'Box', NULL, 3, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 23:33:13');
INSERT INTO `menu` VALUES (4, '采购管理', 'Procurement', 'PROCUREMENT', NULL, NULL, 'ShoppingCart', NULL, 2, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 23:33:10');
INSERT INTO `menu` VALUES (5, '发货管理', 'Shipping', 'SHIPPING', NULL, NULL, 'Van', NULL, 4, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 23:33:17');
INSERT INTO `menu` VALUES (6, '财务管理', 'Finance', 'FINANCE', NULL, NULL, 'Wallet', NULL, 5, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 23:33:20');
INSERT INTO `menu` VALUES (7, '包材管理', 'Packaging', 'PACKAGING', NULL, NULL, 'Box', NULL, 6, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 23:33:30');
INSERT INTO `menu` VALUES (8, '用户管理', 'Users', 'SYSTEM_USERS', NULL, '', 'User', NULL, 7, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-19 01:13:02');
INSERT INTO `menu` VALUES (9, '审计日志', 'Audit Logs', 'SYSTEM_AUDIT', 1, '/system/audit-logs', NULL, NULL, 120, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 20:17:48');
INSERT INTO `menu` VALUES (10, '系统配置', 'Settings', 'SYSTEM_SETTINGS', 1, '/system/settings', NULL, NULL, 130, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 20:17:48');
INSERT INTO `menu` VALUES (11, '字段标签', 'Field Labels', 'SYSTEM_FIELD_LABELS', 1, '/system/field-labels', NULL, NULL, 140, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 20:17:48');
INSERT INTO `menu` VALUES (12, '产品列表', 'Products', 'PRODUCT_SKUS', 2, '/product/list', NULL, NULL, 1, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-20 18:04:20');
INSERT INTO `menu` VALUES (14, '父体管理', 'Parent Products', 'PRODUCT_PARENTS', 2, '/product/parents', NULL, NULL, 230, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 22:23:59');
INSERT INTO `menu` VALUES (15, '仓库管理', 'Warehouses', 'INVENTORY_WAREHOUSES', 3, '/inventory/warehouses', NULL, NULL, 400, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-22 16:38:01');
INSERT INTO `menu` VALUES (16, '库存余额', 'Balances', 'INVENTORY_BALANCES', 3, '/inventory/balances', NULL, NULL, 320, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 20:17:48');
INSERT INTO `menu` VALUES (17, '库存流水', 'Movements', 'INVENTORY_MOVEMENTS', 3, '/inventory/movements', NULL, NULL, 330, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 20:17:48');
INSERT INTO `menu` VALUES (18, '新增流水', 'New Movement', 'INVENTORY_MOVEMENTS_CREATE', 3, '/inventory/movements/create', NULL, NULL, 340, 1, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 20:17:48');
INSERT INTO `menu` VALUES (19, '采购订单', 'Purchase Orders', 'PROCUREMENT_ORDERS', 4, '/procurement/purchase-orders', NULL, NULL, 1, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-22 16:42:26');
INSERT INTO `menu` VALUES (20, '发货单', 'Shipments', 'SHIPPING_SHIPMENTS', 5, '/shipping/shipments', NULL, NULL, 0, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-25 22:28:30');
INSERT INTO `menu` VALUES (21, '现金流水', 'Cash Ledger', 'FINANCE_CASH', 6, '/finance/cash-ledger', NULL, NULL, 610, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 20:17:48');
INSERT INTO `menu` VALUES (22, '成本快照', 'Costing', 'FINANCE_COSTING', 6, '/finance/costing', NULL, NULL, 620, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 20:17:48');
INSERT INTO `menu` VALUES (23, '包材物料', 'Items', 'PACKAGING_ITEMS', 7, '/packaging/items', NULL, NULL, 710, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 20:17:48');
INSERT INTO `menu` VALUES (24, '包材流水', 'Ledger', 'PACKAGING_LEDGER', 7, '/packaging/ledger', NULL, NULL, 720, 0, NULL, 'ACTIVE', '2026-01-18 20:13:49', '2026-01-18 20:17:48');
INSERT INTO `menu` VALUES (25, '组合产品', 'Product Combos', 'PRODUCT_COMBOS', 2, '/product/combos', 'Collection', 'modules/product/views/ProductComboList.vue', 40, 0, 'product.combo', 'ACTIVE', '2026-01-18 22:19:10', '2026-01-18 22:19:10');
INSERT INTO `menu` VALUES (26, '菜单列表', 'Menu List', 'system-menus', 1, '/system/menus', 'Menu', 'System/MenuList', 20, 0, 'system:menu:list', 'ACTIVE', '2026-01-19 21:49:40', '2026-01-19 21:49:40');
INSERT INTO `menu` VALUES (27, '供应商管理', 'Supplier', 'supplier', NULL, '/supplier', 'Box', '', 80, 0, 'supplier:manage', 'ACTIVE', '2026-01-20 17:50:42', '2026-01-20 17:52:57');
INSERT INTO `menu` VALUES (28, '供应商列表', 'Suppliers', 'supplier-list', 27, '/supplier/suppliers', 'User', 'supplier/SupplierList', 81, 0, 'supplier:read', 'ACTIVE', '2026-01-20 17:50:42', '2026-01-20 17:50:42');
INSERT INTO `menu` VALUES (29, '供应商报价', 'Supplier Quotes', 'supplier-quotes', 27, '/supplier/product-quotes', 'Document', 'supplier/SupplierProductQuoteList', 82, 0, 'supplier:quote', 'ACTIVE', '2026-01-20 17:50:42', '2026-01-20 17:50:42');
INSERT INTO `menu` VALUES (30, '打包管理', 'Assembly Management', 'procurement:assembly', 4, '/procurement/assembly', 'Tools', NULL, 2, 0, NULL, 'ACTIVE', '2026-01-22 16:40:57', '2026-01-22 16:42:30');
INSERT INTO `menu` VALUES (31, '装箱规格', 'Package Specs', 'shipping.package-specs', 5, '/shipping/package-specs', 'Box', 'shipping/views/PackageSpecList', 1, 0, '', 'ACTIVE', '2026-01-25 22:27:57', '2026-01-25 22:28:32');
INSERT INTO `menu` VALUES (32, '物流管理', 'logistics', 'Truck', NULL, '', 'Location', '', 99, 0, '', 'ACTIVE', '2026-01-27 14:03:36', '2026-01-27 16:13:33');
INSERT INTO `menu` VALUES (33, '物流供应商', NULL, 'logistics_providers', 32, '/logistics/providers', 'el-icon-shop', 'logistics/providers/index', 1, 0, NULL, 'ACTIVE', '2026-01-27 14:06:18', '2026-01-27 19:07:24');
INSERT INTO `menu` VALUES (34, '运费报价', NULL, 'logistics_shipping_rates', 32, '/logistics/shipping-rates', 'el-icon-money', 'logistics/shipping-rates/index', 3, 0, NULL, 'ACTIVE', '2026-01-27 14:06:18', '2026-01-27 19:06:59');
INSERT INTO `menu` VALUES (35, '物流服务', NULL, 'logistics_services', 32, '/logistics/services', NULL, NULL, 2, 0, NULL, 'ACTIVE', '2026-01-27 19:00:12', '2026-01-27 19:07:25');

-- ----------------------------
-- Table structure for package_spec
-- ----------------------------
DROP TABLE IF EXISTS `package_spec`;
CREATE TABLE `package_spec`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `length` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '长(cm)',
  `width` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '宽(cm)',
  `height` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '高(cm)',
  `weight` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '重量(kg)',
  `quantity_per_box` int UNSIGNED NOT NULL DEFAULT 1 COMMENT '每箱产品数\r\n        量',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '备注',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
  `created_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '更新人',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_gmt_modified`(`gmt_modified` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '装箱规格' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of package_spec
-- ----------------------------
INSERT INTO `package_spec` VALUES (1, '组合1号', 50.00, 40.00, 38.00, 22.00, 1, NULL, 'ACTIVE', NULL, NULL, '2026-01-25 22:29:25', '2026-01-26 19:24:04');

-- ----------------------------
-- Table structure for package_spec_packaging_items
-- ----------------------------
DROP TABLE IF EXISTS `package_spec_packaging_items`;
CREATE TABLE `package_spec_packaging_items`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `package_spec_id` bigint UNSIGNED NOT NULL COMMENT '装箱规格ID (关联 package_spec 表)',
  `packaging_item_id` bigint UNSIGNED NOT NULL COMMENT '包材ID (关联 packaging_item 表)',
  `quantity_per_box` decimal(10, 3) NOT NULL COMMENT '每箱需要的包材数量',
  `notes` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注说明',
  `created_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '更新人',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_package_spec_packaging`(`package_spec_id` ASC, `packaging_item_id` ASC) USING BTREE,
  INDEX `idx_package_spec_id`(`package_spec_id` ASC) USING BTREE,
  INDEX `idx_packaging_item_id`(`packaging_item_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '装箱规格包材关联表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of package_spec_packaging_items
-- ----------------------------

-- ----------------------------
-- Table structure for packaging_item
-- ----------------------------
DROP TABLE IF EXISTS `packaging_item`;
CREATE TABLE `packaging_item`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `trace_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '追踪ID',
  `item_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '物料编码',
  `item_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '物料名称',
  `category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '类别: BOX(纸箱), BAG(包装袋), TAPE(胶带), LABEL(标签), BUBBLE_WRAP(气泡膜), FILLER(填充物), OTHER(其他)',
  `specification` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '规格描述',
  `unit_cost` decimal(10, 4) NOT NULL DEFAULT 0.0000 COMMENT '单位成本',
  `currency` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'CNY' COMMENT '货币',
  `unit` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'PCS' COMMENT '单位: PCS(个), ROLL(卷), METER(米), KG(千克)等',
  `quantity_on_hand` bigint UNSIGNED NOT NULL DEFAULT 0 COMMENT '库存数量',
  `reorder_point` bigint UNSIGNED NULL DEFAULT NULL COMMENT '补货点',
  `reorder_quantity` bigint UNSIGNED NULL DEFAULT NULL COMMENT '建议补货数量',
  `supplier_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '供应商名称',
  `supplier_contact` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '供应商联系方式',
  `status` enum('ACTIVE','INACTIVE') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '备注',
  `created_by` bigint UNSIGNED NOT NULL COMMENT '创建人',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_item_code`(`item_code` ASC) USING BTREE,
  INDEX `idx_trace_id`(`trace_id` ASC) USING BTREE,
  INDEX `idx_category`(`category` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_item_name`(`item_name` ASC) USING BTREE,
  INDEX `idx_created_by`(`created_by` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '包装材料主数据表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of packaging_item
-- ----------------------------
INSERT INTO `packaging_item` VALUES (1, 'PKG-1769674114-8', 'BOX-001', '纸箱', 'BOX', '20x15x10cm', 2.5000, 'CNY', 'PCS', 100, NULL, NULL, 'test', NULL, 'ACTIVE', NULL, 8, '2026-01-29 16:08:34', '2026-01-29 16:15:59');
INSERT INTO `packaging_item` VALUES (4, 'PKG-1769674567-8', 'BOX-002', 'FBA外纸箱', 'BOX', NULL, 90000.0000, 'CNY', 'PCS', 10, NULL, NULL, '广州电子元件', NULL, 'ACTIVE', NULL, 8, '2026-01-29 16:16:08', '2026-01-29 16:25:12');

-- ----------------------------
-- Table structure for packaging_ledger
-- ----------------------------
DROP TABLE IF EXISTS `packaging_ledger`;
CREATE TABLE `packaging_ledger`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `trace_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '追踪ID',
  `packaging_item_id` bigint UNSIGNED NOT NULL COMMENT '包装物料ID',
  `transaction_type` enum('IN','OUT','ADJUSTMENT') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '流水类型: IN-入库, OUT-出库, ADJUSTMENT-调整',
  `quantity` bigint NOT NULL COMMENT '数量（正数表示入库，负数表示出库）',
  `unit_cost` decimal(10, 4) NOT NULL DEFAULT 0.0000 COMMENT '单位成本',
  `total_cost` decimal(15, 2) GENERATED ALWAYS AS ((abs(`quantity`) * `unit_cost`)) STORED COMMENT '总成本（自动计算）' NULL,
  `quantity_before` bigint UNSIGNED NOT NULL COMMENT '操作前库存',
  `quantity_after` bigint UNSIGNED NOT NULL COMMENT '操作后库存',
  `reference_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '关联单据类型: SHIPMENT(发货单), PURCHASE_ORDER(采购单), MANUAL(手工调整)等',
  `reference_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '关联单据ID',
  `occurred_at` datetime NOT NULL COMMENT '发生日期',
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '备注',
  `created_by` bigint UNSIGNED NOT NULL COMMENT '创建人',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_trace_id`(`trace_id` ASC) USING BTREE,
  INDEX `idx_packaging_item_id`(`packaging_item_id` ASC) USING BTREE,
  INDEX `idx_transaction_type`(`transaction_type` ASC) USING BTREE,
  INDEX `idx_occurred_at`(`occurred_at` ASC) USING BTREE,
  INDEX `idx_reference_type_reference_id`(`reference_type` ASC, `reference_id` ASC) USING BTREE,
  INDEX `idx_created_by`(`created_by` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '包装材料流水表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of packaging_ledger
-- ----------------------------
INSERT INTO `packaging_ledger` VALUES (1, 'PKG-IN-1769674274-8', 1, 'IN', 100, 2.5000, DEFAULT, 0, 100, NULL, NULL, '2026-01-29 16:11:15', NULL, 8, '2026-01-29 16:11:15', '2026-01-29 16:11:14');
INSERT INTO `packaging_ledger` VALUES (2, 'PKG-IN-1769675112-8', 4, 'IN', 10, 9.0000, DEFAULT, 0, 10, NULL, NULL, '2026-01-29 16:25:13', NULL, 8, '2026-01-29 16:25:13', '2026-01-29 16:25:12');

-- ----------------------------
-- Table structure for permission
-- ----------------------------
DROP TABLE IF EXISTS `permission`;
CREATE TABLE `permission`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限名称',
  `code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限编码',
  `module` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模块',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '描述',
  `status` enum('ACTIVE','DISABLED') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_code`(`code` ASC) USING BTREE,
  INDEX `idx_module`(`module` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '权限表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of permission
-- ----------------------------

-- ----------------------------
-- Table structure for product
-- ----------------------------
DROP TABLE IF EXISTS `product`;
CREATE TABLE `product`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'SKU ID',
  `seller_sku` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '卖家SKU',
  `asin` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'ASIN',
  `title` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品标题',
  `fnsku` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'FNSKU',
  `marketplace` enum('US','CA','AU','UK','DE','JP') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '站点',
  `parent_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '父体ID',
  `combo_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '组合ID',
  `is_combo_main` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '主组合产品标记（1是，0否）',
  `supplier_id` bigint UNSIGNED NOT NULL COMMENT '默认供应商ID',
  `unit_cost` decimal(15, 4) NULL DEFAULT NULL COMMENT '单位成本',
  `weight` decimal(10, 2) NULL DEFAULT NULL COMMENT '重量(kg)',
  `dimensions` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '尺寸(LxWxH cm)',
  `status` enum('ACTIVE','INACTIVE','DISCONTINUED') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
  `image_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '主图链接',
  `images` json NULL COMMENT '产品图片数组（JSON格式，最多10张）',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '备注',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_seller_sku_marketplace`(`seller_sku` ASC, `marketplace` ASC) USING BTREE,
  INDEX `idx_asin`(`asin` ASC) USING BTREE,
  INDEX `idx_marketplace`(`marketplace` ASC) USING BTREE,
  INDEX `idx_supplier_id`(`supplier_id` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_parent_id`(`parent_id` ASC) USING BTREE,
  INDEX `idx_combo_id`(`combo_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'SKU表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of product
-- ----------------------------
INSERT INTO `product` VALUES (1, 'SKU-MAIN-001', 'ASIN-MAIN-001', '组合主产品', 'FNSKU-M-001', 'US', 1, 1, 1, 1, 19.9900, 0.50, '10x10x10', 'ACTIVE', '', '[]', '组合主产品', '2026-01-18 23:41:07', '2026-01-21 15:07:02');
INSERT INTO `product` VALUES (2, 'SKU-CHILD-002', 'ASIN-CHILD-002', '组合子产品A', 'FNSKU-C-002', 'US', 1, 1, 0, 1, 9.9900, 0.30, '8x8x8', 'ACTIVE', '', '[]', '组合子产品A', '2026-01-18 23:41:07', '2026-01-21 15:07:02');
INSERT INTO `product` VALUES (3, 'SKU-CHILD-003', 'ASIN-CHILD-003', '组合子产品B', 'FNSKU-C-003', 'US', 1, 1, 0, 1, 12.9900, 0.35, '9x9x9', 'ACTIVE', '', '[]', '组合子产品B', '2026-01-18 23:41:07', '2026-01-21 15:07:02');
INSERT INTO `product` VALUES (4, 'SKU-TEST-101', 'ASIN-TEST-101', 'Test Product 101', '', 'US', NULL, NULL, 0, 1, 9.9900, 0.50, '10x8x2', 'ACTIVE', '/uploads/products/20260120013259-新对话.png', '[]', 'seed', '2026-01-19 01:06:15', '2026-01-20 17:13:21');
INSERT INTO `product` VALUES (5, 'SKU-TEST-102', 'ASIN-TEST-102', 'Test Product 102', '', 'CA', NULL, NULL, 0, 1, 12.5000, 0.60, '12x9x3', 'ACTIVE', '', '[]', 'seed', '2026-01-19 01:06:15', '2026-01-20 17:13:22');
INSERT INTO `product` VALUES (6, 'SKU-TEST-103', 'ASIN-TEST-103', 'Test Product 103', '', 'UK', NULL, NULL, 0, 3, 7.3000, 0.40, '9x7x2', 'ACTIVE', '', '[]', 'seed', '2026-01-19 01:06:15', '2026-01-20 18:35:41');

-- ----------------------------
-- Table structure for product_combo
-- ----------------------------
DROP TABLE IF EXISTS `product_combo`;
CREATE TABLE `product_combo`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '组合记录ID',
  `combo_id` bigint UNSIGNED NOT NULL COMMENT '组合ID',
  `main_product_id` bigint UNSIGNED NOT NULL COMMENT '主产品ID',
  `product_id` bigint UNSIGNED NOT NULL COMMENT '组件产品ID',
  `qty_ratio` int UNSIGNED NOT NULL DEFAULT 1 COMMENT '数量比例',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_combo_product`(`combo_id` ASC, `product_id` ASC) USING BTREE,
  INDEX `idx_combo_id`(`combo_id` ASC) USING BTREE,
  INDEX `idx_main_product_id`(`main_product_id` ASC) USING BTREE,
  INDEX `idx_product_id`(`product_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '产品组合表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of product_combo
-- ----------------------------
INSERT INTO `product_combo` VALUES (3, 1, 1, 1, 1, '2026-01-21 15:07:02', '2026-01-21 15:07:02');
INSERT INTO `product_combo` VALUES (4, 1, 1, 2, 1, '2026-01-21 15:07:02', '2026-01-21 15:07:02');
INSERT INTO `product_combo` VALUES (5, 1, 1, 3, 1, '2026-01-21 15:07:02', '2026-01-21 15:07:02');

-- ----------------------------
-- Table structure for product_image
-- ----------------------------
DROP TABLE IF EXISTS `product_image`;
CREATE TABLE `product_image`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `product_id` bigint UNSIGNED NOT NULL COMMENT '产品ID',
  `image_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '图片URL',
  `sort_order` int UNSIGNED NOT NULL DEFAULT 1 COMMENT '排序序号',
  `is_primary` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否主图(1是/0否)',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '状态',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '备注',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_product_image`(`product_id` ASC, `image_url` ASC) USING BTREE,
  INDEX `idx_product_id_sort`(`product_id` ASC, `sort_order` ASC) USING BTREE,
  INDEX `idx_product_id_primary`(`product_id` ASC, `is_primary` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 18 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '产品图片表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of product_image
-- ----------------------------
INSERT INTO `product_image` VALUES (12, 4, '/uploads/products/20260120013259-新对话.png', 1, 1, '', '', '2026-01-20 01:38:55', '2026-01-20 01:38:55');
INSERT INTO `product_image` VALUES (13, 4, '/uploads/products/20260120013259-1.png', 2, 0, '', '', '2026-01-20 01:38:55', '2026-01-20 01:38:55');
INSERT INTO `product_image` VALUES (14, 4, '/uploads/products/20260120013259-2.png', 3, 0, '', '', '2026-01-20 01:38:55', '2026-01-20 01:38:55');
INSERT INTO `product_image` VALUES (15, 4, '/uploads/products/20260120013259-3.png', 4, 0, '', '', '2026-01-20 01:38:55', '2026-01-20 01:38:55');
INSERT INTO `product_image` VALUES (16, 4, '/uploads/products/20260120013259-4.png', 5, 0, '', '', '2026-01-20 01:38:55', '2026-01-20 01:38:55');
INSERT INTO `product_image` VALUES (17, 4, '/uploads/products/20260120013852-新对话.png', 6, 0, '', '', '2026-01-20 01:38:55', '2026-01-20 01:38:55');

-- ----------------------------
-- Table structure for product_packaging_items
-- ----------------------------
DROP TABLE IF EXISTS `product_packaging_items`;
CREATE TABLE `product_packaging_items`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `product_id` bigint UNSIGNED NOT NULL COMMENT '产品ID (关联 sku 表)',
  `packaging_item_id` bigint UNSIGNED NOT NULL COMMENT '包材ID (关联 packaging_item 表)',
  `quantity_per_unit` decimal(10, 3) NOT NULL COMMENT '每个产品需要的包材数量',
  `notes` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注说明',
  `created_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '更新人',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_product_packaging`(`product_id` ASC, `packaging_item_id` ASC) USING BTREE,
  INDEX `idx_product_id`(`product_id` ASC) USING BTREE,
  INDEX `idx_packaging_item_id`(`packaging_item_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '产品包材关联表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of product_packaging_items
-- ----------------------------

-- ----------------------------
-- Table structure for product_parent
-- ----------------------------
DROP TABLE IF EXISTS `product_parent`;
CREATE TABLE `product_parent`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '父体ID',
  `parent_asin` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Parent ASIN',
  `title` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品标题',
  `marketplace` enum('US','CA','AU','UK','DE','JP') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '站点',
  `brand` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '品牌',
  `category` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '类目',
  `status` enum('ACTIVE','INACTIVE','DISCONTINUED') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
  `image_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '主图链接',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '备注',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_parent_asin_marketplace`(`parent_asin` ASC, `marketplace` ASC) USING BTREE,
  INDEX `idx_marketplace`(`marketplace` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'Listing表（父体）' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of product_parent
-- ----------------------------
INSERT INTO `product_parent` VALUES (1, 'PARENT-ASIN-001', '父体测试商品', 'US', 'TestBrand', 'TestCategory', 'ACTIVE', '', '用于父体功能测试', '2026-01-18 23:41:07', '2026-01-18 23:41:07');

-- ----------------------------
-- Table structure for product_supplier
-- ----------------------------
DROP TABLE IF EXISTS `product_supplier`;
CREATE TABLE `product_supplier`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `product_id` bigint UNSIGNED NOT NULL COMMENT '产品ID',
  `supplier_id` bigint UNSIGNED NOT NULL COMMENT '供应商ID',
  `is_default` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否默认供应商',
  `purchase_price` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '供应价',
  `payment_terms` int UNSIGNED NULL DEFAULT NULL COMMENT '账期(天)',
  `moq` decimal(10, 2) NULL DEFAULT NULL COMMENT '最小起订量',
  `lead_time_days` int UNSIGNED NULL DEFAULT NULL COMMENT '交期(天)',
  `currency` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '币种',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_product_supplier`(`product_id` ASC, `supplier_id` ASC) USING BTREE,
  INDEX `idx_product_supplier_product`(`product_id` ASC) USING BTREE,
  INDEX `idx_product_supplier_supplier`(`supplier_id` ASC) USING BTREE,
  INDEX `idx_product_supplier_default`(`product_id` ASC, `is_default` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '产品-供应商关系' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of product_supplier
-- ----------------------------
INSERT INTO `product_supplier` VALUES (1, 1, 1, 1, 12.50, 30, 100.00, 15, 'USD', 'ACTIVE', '2026-01-19 16:25:59', '2026-01-19 16:25:59');
INSERT INTO `product_supplier` VALUES (2, 1, 4, 0, 11.80, 20, 200.00, 20, 'USD', 'ACTIVE', '2026-01-19 16:25:59', '2026-01-19 16:25:59');
INSERT INTO `product_supplier` VALUES (3, 2, 1, 1, 8.90, 30, 50.00, 12, 'USD', 'ACTIVE', '2026-01-19 16:25:59', '2026-01-19 16:25:59');
INSERT INTO `product_supplier` VALUES (4, 3, 2, 1, 0.35, 15, 1000.00, 7, 'CNY', 'ACTIVE', '2026-01-19 16:25:59', '2026-01-19 16:25:59');

-- ----------------------------
-- Table structure for product_supplier_quote
-- ----------------------------
DROP TABLE IF EXISTS `product_supplier_quote`;
CREATE TABLE `product_supplier_quote`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '报价ID',
  `product_id` bigint UNSIGNED NOT NULL COMMENT '产品ID',
  `supplier_id` bigint UNSIGNED NOT NULL COMMENT '供应商ID',
  `price` decimal(15, 4) NOT NULL COMMENT '报价',
  `currency` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '币种',
  `qty_moq` int UNSIGNED NOT NULL DEFAULT 1 COMMENT '起订量',
  `lead_time_days` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '交期(天)',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '备注',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_product_supplier`(`product_id` ASC, `supplier_id` ASC) USING BTREE,
  INDEX `idx_product_id`(`product_id` ASC) USING BTREE,
  INDEX `idx_supplier_id`(`supplier_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '产品供应商报价表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of product_supplier_quote
-- ----------------------------
INSERT INTO `product_supplier_quote` VALUES (1, 6, 2, 1.0000, 'USD', 1, 0, 'ACTIVE', '', '2026-01-20 18:27:31', '2026-01-20 18:27:31');
INSERT INTO `product_supplier_quote` VALUES (2, 6, 3, 3.0000, 'USD', 1, 0, 'ACTIVE', '', '2026-01-20 18:27:43', '2026-01-20 18:27:43');
INSERT INTO `product_supplier_quote` VALUES (3, 1, 1, 100.0000, 'USD', 100, 0, 'ACTIVE', '', '2026-01-21 15:09:41', '2026-01-21 18:30:34');
INSERT INTO `product_supplier_quote` VALUES (4, 2, 1, 1.0000, 'USD', 1, 0, 'ACTIVE', '', '2026-01-21 15:09:46', '2026-01-21 18:37:41');
INSERT INTO `product_supplier_quote` VALUES (5, 3, 2, 1.0000, 'USD', 1, 0, 'ACTIVE', '', '2026-01-21 15:09:49', '2026-01-21 18:37:47');
INSERT INTO `product_supplier_quote` VALUES (6, 5, 3, 3.0000, 'USD', 1, 0, 'ACTIVE', '', '2026-01-21 15:09:53', '2026-01-21 18:37:53');
INSERT INTO `product_supplier_quote` VALUES (7, 4, 4, 2.0000, 'USD', 1, 0, 'ACTIVE', '', '2026-01-21 15:09:57', '2026-01-21 18:37:59');
INSERT INTO `product_supplier_quote` VALUES (8, 1, 2, 10.0000, 'USD', 10, 0, 'ACTIVE', '', '2026-01-21 18:21:24', '2026-01-21 18:30:22');
INSERT INTO `product_supplier_quote` VALUES (9, 1, 4, 1.0000, 'USD', 1, 0, 'ACTIVE', '', '2026-01-21 18:21:30', '2026-01-21 18:30:12');

-- ----------------------------
-- Table structure for purchase_order
-- ----------------------------
DROP TABLE IF EXISTS `purchase_order`;
CREATE TABLE `purchase_order`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '采购单ID',
  `po_number` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '采购单号',
  `supplier_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '供应商ID',
  `marketplace` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '站点',
  `status` enum('DRAFT','ORDERED','SHIPPED','RECEIVED','CLOSED') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'DRAFT' COMMENT '采购单状态',
  `currency` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'USD' COMMENT '币种',
  `total_amount` decimal(12, 4) NOT NULL DEFAULT 0.0000 COMMENT '采购单总金额',
  `ordered_at` datetime NULL DEFAULT NULL COMMENT '下单时间',
  `shipped_at` datetime NULL DEFAULT NULL COMMENT '发货时间',
  `received_at` datetime NULL DEFAULT NULL COMMENT '到货时间',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '备注',
  `created_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '更新人',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_po_number`(`po_number` ASC) USING BTREE,
  INDEX `idx_supplier_id`(`supplier_id` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_marketplace`(`marketplace` ASC) USING BTREE,
  INDEX `idx_gmt_create`(`gmt_create` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 82 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '采购单表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of purchase_order
-- ----------------------------
INSERT INTO `purchase_order` VALUES (29, 'PO202601211905422000', 3, '', 'SHIPPED', 'USD', 6.0000, '2026-01-21 19:07:43', '2026-01-21 19:07:50', NULL, '', 8, 8, '2026-01-21 19:05:43', '2026-01-21 19:07:50');
INSERT INTO `purchase_order` VALUES (40, 'PO202601212000498500', 2, 'US', 'RECEIVED', 'USD', 1.0000, '2026-01-21 20:01:33', '2026-01-21 20:01:36', '2026-01-21 20:04:13', '', 8, 8, '2026-01-21 20:00:50', '2026-01-21 20:04:13');
INSERT INTO `purchase_order` VALUES (44, 'PO202601221716443200', 1, '', 'RECEIVED', 'USD', 10001.0000, '2026-01-22 17:17:12', '2026-01-22 17:17:14', '2026-01-22 17:17:29', '', 8, 8, '2026-01-22 17:16:45', '2026-01-22 17:17:29');
INSERT INTO `purchase_order` VALUES (78, 'PO20260123141609864d37d4', 1, 'CA', 'RECEIVED', 'USD', 10001.0000, '2026-01-24 12:27:03', '2026-01-24 12:27:08', '2026-01-24 12:27:13', '', 8, 8, '2026-01-23 14:16:09', '2026-01-24 12:27:13');
INSERT INTO `purchase_order` VALUES (79, 'PO202601231416094da24355', 2, 'CA', 'RECEIVED', 'USD', 1.0000, '2026-01-24 12:26:50', '2026-01-24 12:26:56', '2026-01-24 12:27:01', '', 8, 8, '2026-01-23 14:16:09', '2026-01-24 12:27:01');
INSERT INTO `purchase_order` VALUES (80, 'PO20260123141628c78db3ce', 2, '', 'RECEIVED', 'USD', 1.0000, '2026-01-23 14:17:34', '2026-01-23 14:17:36', '2026-01-23 14:17:44', '', 8, 8, '2026-01-23 14:16:29', '2026-01-23 14:17:44');
INSERT INTO `purchase_order` VALUES (81, 'PO202601231416282515eabc', 1, '', 'RECEIVED', 'USD', 10001.0000, '2026-01-23 14:17:54', '2026-01-23 14:17:57', '2026-01-23 14:18:09', '', 8, 8, '2026-01-23 14:16:29', '2026-01-23 14:18:09');

-- ----------------------------
-- Table structure for purchase_order_item
-- ----------------------------
DROP TABLE IF EXISTS `purchase_order_item`;
CREATE TABLE `purchase_order_item`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '采购单明细ID',
  `purchase_order_id` bigint UNSIGNED NOT NULL COMMENT '采购单ID',
  `sku_id` bigint UNSIGNED NOT NULL COMMENT 'SKU ID',
  `qty_ordered` int UNSIGNED NOT NULL COMMENT '采购数量',
  `qty_received` int UNSIGNED NOT NULL COMMENT '到货数量',
  `unit_cost` decimal(12, 4) NOT NULL DEFAULT 0.0000 COMMENT '单价',
  `currency` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'USD' COMMENT '币种',
  `subtotal` decimal(12, 4) NOT NULL DEFAULT 0.0000 COMMENT '小计金额',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_purchase_order_id`(`purchase_order_id` ASC) USING BTREE,
  INDEX `idx_sku_id`(`sku_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 74 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '采购单明细表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of purchase_order_item
-- ----------------------------
INSERT INTO `purchase_order_item` VALUES (19, 29, 5, 1, 0, 3.0000, 'USD', 3.0000, '2026-01-21 19:05:43', '2026-01-21 19:05:43');
INSERT INTO `purchase_order_item` VALUES (20, 29, 6, 1, 0, 3.0000, 'USD', 3.0000, '2026-01-21 19:05:43', '2026-01-21 19:05:43');
INSERT INTO `purchase_order_item` VALUES (28, 40, 3, 1, 1, 1.0000, 'USD', 1.0000, '2026-01-21 20:00:50', '2026-01-21 20:00:50');
INSERT INTO `purchase_order_item` VALUES (33, 44, 1, 100, 100, 100.0000, 'USD', 10000.0000, '2026-01-22 17:16:45', '2026-01-22 17:16:45');
INSERT INTO `purchase_order_item` VALUES (34, 44, 2, 1, 1, 1.0000, 'USD', 1.0000, '2026-01-22 17:16:45', '2026-01-22 17:16:45');
INSERT INTO `purchase_order_item` VALUES (68, 78, 1, 100, 100, 100.0000, 'USD', 10000.0000, '2026-01-23 14:16:09', '2026-01-23 14:16:09');
INSERT INTO `purchase_order_item` VALUES (69, 78, 2, 1, 1, 1.0000, 'USD', 1.0000, '2026-01-23 14:16:09', '2026-01-23 14:16:09');
INSERT INTO `purchase_order_item` VALUES (70, 79, 3, 1, 1, 1.0000, 'USD', 1.0000, '2026-01-23 14:16:09', '2026-01-23 14:16:09');
INSERT INTO `purchase_order_item` VALUES (71, 80, 3, 1, 1, 1.0000, 'USD', 1.0000, '2026-01-23 14:16:29', '2026-01-23 14:16:29');
INSERT INTO `purchase_order_item` VALUES (72, 81, 1, 100, 100, 100.0000, 'USD', 10000.0000, '2026-01-23 14:16:29', '2026-01-23 14:16:29');
INSERT INTO `purchase_order_item` VALUES (73, 81, 2, 1, 1, 1.0000, 'USD', 1.0000, '2026-01-23 14:16:29', '2026-01-23 14:16:29');

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称（ADMIN/OPERATOR/VIEWER）',
  `display_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '显示名称',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '角色说明',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_name`(`name` ASC) USING BTREE,
  INDEX `idx_gmt_create`(`gmt_create` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 17 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of role
-- ----------------------------
INSERT INTO `role` VALUES (16, 'admin', '超级管理员', '系统超级管理员，拥有所有权限', '2026-01-19 19:35:19', '2026-01-19 19:35:19');

-- ----------------------------
-- Table structure for role_permission
-- ----------------------------
DROP TABLE IF EXISTS `role_permission`;
CREATE TABLE `role_permission`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `role_id` bigint UNSIGNED NOT NULL COMMENT '角色ID',
  `permission_id` bigint UNSIGNED NOT NULL COMMENT '权限ID',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_role_id_permission_id`(`role_id` ASC, `permission_id` ASC) USING BTREE,
  INDEX `idx_role_id`(`role_id` ASC) USING BTREE,
  INDEX `idx_permission_id`(`permission_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色权限关联表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of role_permission
-- ----------------------------

-- ----------------------------
-- Table structure for setting
-- ----------------------------
DROP TABLE IF EXISTS `setting`;
CREATE TABLE `setting`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `scope_type` enum('GLOBAL','SITE','LISTING','SKU') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'GLOBAL' COMMENT '作用域类型',
  `scope_id` bigint UNSIGNED NOT NULL DEFAULT 0 COMMENT '作用域ID（GLOBAL=0）',
  `setting_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置Key',
  `setting_type` enum('SYSTEM','CUSTOM') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'CUSTOM' COMMENT '配置类型',
  `value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置值（字符串）',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '配置说明',
  `created_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '更新人',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_scope_key`(`scope_type` ASC, `scope_id` ASC, `setting_key` ASC) USING BTREE,
  INDEX `idx_setting_key`(`setting_key` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统配置表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of setting
-- ----------------------------

-- ----------------------------
-- Table structure for shipment
-- ----------------------------
DROP TABLE IF EXISTS `shipment`;
CREATE TABLE `shipment`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `shipment_number` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '发货单号',
  `order_number` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '订单号/Reference',
  `sales_channel` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '销售渠道: Amazon/eBay/Shopify/Offline等',
  `warehouse_id` bigint UNSIGNED NOT NULL COMMENT '发货仓库ID',
  `destination_warehouse_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '目的地仓库ID',
  `destination_type` enum('PLATFORM_WAREHOUSE','CUSTOMER','OWN_WAREHOUSE','SUPPLIER','OTHER') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'PLATFORM_WAREHOUSE' COMMENT '收货方类型',
  `destination_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '收货方名称',
  `destination_contact` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '收货联系人',
  `destination_phone` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '收货电话',
  `destination_address` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '收货地址',
  `logistics_provider_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '物流供应商ID',
  `shipping_rate_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '运费报价ID',
  `transport_mode` enum('EXPRESS','AIR','SEA','RAIL','TRUCK') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '运输方式',
  `destination_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '收货方代码(如FBA仓库代码)',
  `carrier` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '承运商',
  `shipping_method` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '运输方式: 快递/空运/海运/陆运',
  `tracking_number` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '物流追踪号(可能多个,逗号分隔)',
  `box_count` int UNSIGNED NULL DEFAULT 0 COMMENT '箱数',
  `total_weight` decimal(10, 2) NULL DEFAULT 0.00 COMMENT '总重量(kg)',
  `total_volume` decimal(10, 3) NULL DEFAULT 0.000 COMMENT '总体积(m³)',
  `shipping_cost` decimal(12, 4) NOT NULL DEFAULT 0.0000 COMMENT '运费',
  `currency` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'USD' COMMENT '币种',
  `ship_date` date NULL DEFAULT NULL COMMENT '发货日期',
  `expected_delivery_date` date NULL DEFAULT NULL COMMENT '预计到达日期',
  `actual_delivery_date` date NULL DEFAULT NULL COMMENT '实际到达日期',
  `status` enum('DRAFT','CONFIRMED','PACKED','SHIPPED','DELIVERED','CANCELLED') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'DRAFT' COMMENT '状态',
  `confirmed_at` datetime NULL DEFAULT NULL COMMENT '确认时间',
  `shipped_at` datetime NULL DEFAULT NULL COMMENT '发货时间',
  `delivered_at` datetime NULL DEFAULT NULL COMMENT '送达时间',
  `confirmed_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '确认人',
  `shipped_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '发货人',
  `delivered_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '签收人',
  `inventory_locked` tinyint(1) NOT NULL DEFAULT 0 COMMENT '库存是否已锁定',
  `inventory_deducted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '库存是否已扣减',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '备注',
  `internal_notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '内部备注',
  `created_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '更新人',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `shipment_number`(`shipment_number` ASC) USING BTREE,
  INDEX `idx_order_number`(`order_number` ASC) USING BTREE,
  INDEX `idx_warehouse_id`(`warehouse_id` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_destination_type`(`destination_type` ASC) USING BTREE,
  INDEX `idx_tracking_number`(`tracking_number` ASC) USING BTREE,
  INDEX `idx_ship_date`(`ship_date` ASC) USING BTREE,
  INDEX `idx_gmt_create`(`gmt_create` ASC) USING BTREE,
  INDEX `idx_destination_warehouse`(`destination_warehouse_id` ASC) USING BTREE,
  INDEX `idx_logistics_provider`(`logistics_provider_id` ASC) USING BTREE,
  INDEX `idx_shipping_rate`(`shipping_rate_id` ASC) USING BTREE,
  INDEX `idx_transport_mode`(`transport_mode` ASC) USING BTREE,
  INDEX `idx_confirmed_at`(`confirmed_at` ASC) USING BTREE,
  INDEX `idx_shipped_at`(`shipped_at` ASC) USING BTREE,
  INDEX `idx_delivered_at`(`delivered_at` ASC) USING BTREE,
  CONSTRAINT `fk_shipment_destination_warehouse` FOREIGN KEY (`destination_warehouse_id`) REFERENCES `warehouse` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_shipment_logistics_provider` FOREIGN KEY (`logistics_provider_id`) REFERENCES `logistics_provider` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_shipment_shipping_rate` FOREIGN KEY (`shipping_rate_id`) REFERENCES `shipping_rate` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '发货单' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of shipment
-- ----------------------------
INSERT INTO `shipment` VALUES (1, 'SH2026012417330212868', NULL, NULL, 1, NULL, 'PLATFORM_WAREHOUSE', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, 0.00, 0.000, 0.0000, 'USD', NULL, NULL, NULL, 'CANCELLED', NULL, NULL, NULL, NULL, NULL, NULL, 0, 0, NULL, NULL, 8, 8, '2026-01-24 17:33:02', '2026-01-26 19:23:24');
INSERT INTO `shipment` VALUES (2, 'SH2026012417544383872', NULL, NULL, 1, NULL, 'PLATFORM_WAREHOUSE', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, 0.00, 0.000, 0.0000, 'USD', NULL, NULL, NULL, 'CANCELLED', NULL, NULL, NULL, NULL, NULL, NULL, 0, 0, NULL, NULL, 8, 8, '2026-01-24 17:54:44', '2026-01-26 19:23:23');
INSERT INTO `shipment` VALUES (3, 'SH2026012515563742913', NULL, NULL, 1, NULL, 'PLATFORM_WAREHOUSE', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, 0.00, 0.000, 0.0000, 'USD', NULL, NULL, NULL, 'CANCELLED', NULL, NULL, NULL, NULL, NULL, NULL, 0, 0, NULL, NULL, 8, 8, '2026-01-25 15:56:37', '2026-01-26 19:23:21');
INSERT INTO `shipment` VALUES (4, 'SH2026012613435899488', NULL, NULL, 1, NULL, 'PLATFORM_WAREHOUSE', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, 0.00, 0.000, 0.0000, 'USD', NULL, NULL, NULL, 'CANCELLED', NULL, NULL, NULL, NULL, NULL, NULL, 0, 0, NULL, NULL, 8, 8, '2026-01-26 13:43:59', '2026-01-26 19:04:43');
INSERT INTO `shipment` VALUES (5, 'SH2026012619245407779', NULL, NULL, 1, NULL, 'PLATFORM_WAREHOUSE', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'test', NULL, 'test', 0, 0.00, 0.000, 0.0000, 'USD', NULL, NULL, NULL, 'DELIVERED', NULL, NULL, NULL, NULL, NULL, NULL, 1, 1, NULL, NULL, 8, 8, '2026-01-26 19:24:54', '2026-01-26 19:25:58');
INSERT INTO `shipment` VALUES (6, 'SH2026012712511026636', NULL, NULL, 1, NULL, 'PLATFORM_WAREHOUSE', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, 0.00, 0.000, 0.0000, 'USD', NULL, NULL, NULL, 'CONFIRMED', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, NULL, NULL, 8, 8, '2026-01-27 12:51:10', '2026-01-27 12:51:25');
INSERT INTO `shipment` VALUES (7, 'SH2026012814372137351', NULL, NULL, 1, NULL, 'PLATFORM_WAREHOUSE', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, 0.00, 0.000, 0.0000, 'USD', NULL, NULL, NULL, 'CONFIRMED', '2026-01-28 15:06:13', NULL, NULL, 8, NULL, NULL, 1, 0, NULL, NULL, 8, 8, '2026-01-28 14:37:21', '2026-01-28 15:06:13');
INSERT INTO `shipment` VALUES (8, 'SH2026012814545113557', NULL, NULL, 1, NULL, 'PLATFORM_WAREHOUSE', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, 0.00, 0.000, 0.0000, 'USD', NULL, NULL, NULL, 'DELIVERED', '2026-01-28 14:55:09', '2026-01-28 14:55:15', '2026-01-28 14:55:19', 8, 8, 8, 1, 1, NULL, NULL, 8, 8, '2026-01-28 14:54:51', '2026-01-28 14:55:19');

-- ----------------------------
-- Table structure for shipment_item
-- ----------------------------
DROP TABLE IF EXISTS `shipment_item`;
CREATE TABLE `shipment_item`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `shipment_id` bigint UNSIGNED NOT NULL COMMENT '发货单ID',
  `sku_id` bigint UNSIGNED NOT NULL COMMENT 'SKU ID',
  `quantity_planned` int UNSIGNED NOT NULL COMMENT '计划发货数量',
  `quantity_shipped` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '实际发货数量',
  `package_spec_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '装箱规格ID',
  `box_quantity` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '装箱数量',
  `unit_cost` decimal(12, 4) NOT NULL DEFAULT 0.0000 COMMENT '单位成本',
  `currency` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'USD' COMMENT '币种',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_shipment_id`(`shipment_id` ASC) USING BTREE,
  INDEX `idx_sku_id`(`sku_id` ASC) USING BTREE,
  INDEX `idx_package_spec_id`(`package_spec_id` ASC) USING BTREE,
  CONSTRAINT `shipment_item_ibfk_1` FOREIGN KEY (`shipment_id`) REFERENCES `shipment` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '发货单明细' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of shipment_item
-- ----------------------------
INSERT INTO `shipment_item` VALUES (1, 1, 1, 0, 0, NULL, 0, 0.0000, 'USD', NULL, '2026-01-24 17:33:02', '2026-01-24 17:33:02');
INSERT INTO `shipment_item` VALUES (2, 2, 2, 1, 0, NULL, 0, 0.0000, 'USD', '', '2026-01-24 17:54:44', '2026-01-24 17:54:44');
INSERT INTO `shipment_item` VALUES (3, 3, 4, 1, 0, NULL, 0, 9.9900, 'USD', '', '2026-01-25 15:56:37', '2026-01-25 15:56:37');
INSERT INTO `shipment_item` VALUES (4, 4, 3, 40, 0, 1, 4, 0.0000, 'USD', '', '2026-01-26 13:43:59', '2026-01-26 13:43:59');
INSERT INTO `shipment_item` VALUES (5, 5, 3, 2, 0, 1, 2, 0.0000, 'USD', '', '2026-01-26 19:24:54', '2026-01-26 19:24:54');
INSERT INTO `shipment_item` VALUES (6, 6, 3, 2, 0, 1, 2, 0.0000, 'USD', '', '2026-01-27 12:51:10', '2026-01-27 12:51:10');
INSERT INTO `shipment_item` VALUES (7, 7, 3, 1, 0, 1, 1, 0.0000, 'USD', '', '2026-01-28 14:37:21', '2026-01-28 14:37:21');
INSERT INTO `shipment_item` VALUES (8, 8, 1, 2, 0, 1, 2, 0.0000, 'USD', '', '2026-01-28 14:54:51', '2026-01-28 14:54:51');

-- ----------------------------
-- Table structure for shipping_rate
-- ----------------------------
DROP TABLE IF EXISTS `shipping_rate`;
CREATE TABLE `shipping_rate`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `service_id` bigint UNSIGNED NOT NULL COMMENT '服务商ID',
  `provider_id` bigint UNSIGNED NOT NULL COMMENT '物流供应商ID',
  `origin_warehouse_id` bigint UNSIGNED NOT NULL COMMENT '起运仓库ID',
  `destination_warehouse_id` bigint UNSIGNED NOT NULL COMMENT '目的仓库ID',
  `transport_mode` enum('EXPRESS','AIR','SEA','RAIL','TRUCK') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '运输方式',
  `pricing_method` enum('PER_KG','PER_CBM','PER_PACKAGE','FIXED') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '计费方式',
  `base_rate` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '基础费率',
  `currency` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'CNY' COMMENT '货币',
  `other_fee` decimal(10, 2) NOT NULL DEFAULT 0.00,
  `min_weight` decimal(10, 2) NULL DEFAULT NULL COMMENT '最小重量(kg)',
  `max_weight` decimal(10, 2) NULL DEFAULT NULL COMMENT '最大重量(kg)',
  `min_volume` decimal(10, 4) NULL DEFAULT NULL COMMENT '最小体积(m³)',
  `max_volume` decimal(10, 4) NULL DEFAULT NULL COMMENT '最大体积(m³)',
  `transit_days` int NULL DEFAULT NULL COMMENT '运输时效(天)',
  `effective_date` date NOT NULL COMMENT '生效日期',
  `expiry_date` date NULL DEFAULT NULL COMMENT '失效日期',
  `status` enum('ACTIVE','INACTIVE','EXPIRED') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '备注',
  `created_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '更新人',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_provider`(`provider_id` ASC) USING BTREE,
  INDEX `idx_origin_warehouse`(`origin_warehouse_id` ASC) USING BTREE,
  INDEX `idx_destination_warehouse`(`destination_warehouse_id` ASC) USING BTREE,
  INDEX `idx_transport_mode`(`transport_mode` ASC) USING BTREE,
  INDEX `idx_effective_date`(`effective_date` ASC) USING BTREE,
  INDEX `idx_expiry_date`(`expiry_date` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_service_id`(`service_id` ASC) USING BTREE,
  CONSTRAINT `shipping_rate_ibfk_1` FOREIGN KEY (`provider_id`) REFERENCES `logistics_provider` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `shipping_rate_ibfk_2` FOREIGN KEY (`origin_warehouse_id`) REFERENCES `warehouse` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `shipping_rate_ibfk_3` FOREIGN KEY (`destination_warehouse_id`) REFERENCES `warehouse` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '运费报价表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of shipping_rate
-- ----------------------------
INSERT INTO `shipping_rate` VALUES (1, 1, 1, 1, 2, 'SEA', 'PER_KG', 8.00, 'USD', 0.00, 50.00, NULL, NULL, NULL, NULL, '2026-01-01', '2026-02-07', 'ACTIVE', '测试', 8, 8, '2026-01-27 20:00:08', '2026-01-27 20:00:08');

-- ----------------------------
-- Table structure for supplier
-- ----------------------------
DROP TABLE IF EXISTS `supplier`;
CREATE TABLE `supplier`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '供应商ID',
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '供应商名称',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ACTIVE' COMMENT '状态(ACTIVE/DISABLED)',
  `supplier_code` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'supplier_code',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '备注',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_supplier_name`(`name` ASC) USING BTREE,
  UNIQUE INDEX `uk_supplier_code`(`supplier_code` ASC) USING BTREE,
  INDEX `idx_supplier_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '供应商表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of supplier
-- ----------------------------
INSERT INTO `supplier` VALUES (1, '深圳优品供应链', 'ACTIVE', '深圳优品供应链', '核心产品供应商', '2026-01-19 16:25:59', '2026-01-20 18:25:27');
INSERT INTO `supplier` VALUES (2, '杭州包装工厂', 'ACTIVE', '杭州包装工厂', '包材合作', '2026-01-19 16:25:59', '2026-01-20 18:25:30');
INSERT INTO `supplier` VALUES (3, '上海物流服务', 'ACTIVE', '上海物流服务', '头程物流', '2026-01-19 16:25:59', '2026-01-20 18:25:33');
INSERT INTO `supplier` VALUES (4, '广州电子元件', 'ACTIVE', '广州电子元件', '电子类供应商', '2026-01-19 16:25:59', '2026-01-20 18:25:37');

-- ----------------------------
-- Table structure for supplier_account
-- ----------------------------
DROP TABLE IF EXISTS `supplier_account`;
CREATE TABLE `supplier_account`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '账户ID',
  `supplier_id` bigint UNSIGNED NOT NULL COMMENT '供应商ID',
  `bank_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '开户行',
  `bank_account` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '银行账号',
  `currency` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '币种',
  `tax_no` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '税号',
  `payment_terms` int UNSIGNED NULL DEFAULT NULL COMMENT '账期(天)',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_supplier_account_supplier`(`supplier_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '供应商账户' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of supplier_account
-- ----------------------------
INSERT INTO `supplier_account` VALUES (1, 1, '中国银行深圳分行', '6222000000000001', 'CNY', '91440300X00000001', 30, '2026-01-19 16:25:59', '2026-01-19 16:25:59');
INSERT INTO `supplier_account` VALUES (2, 2, '建设银行杭州分行', '6227000000000002', 'CNY', '91330100X00000002', 15, '2026-01-19 16:25:59', '2026-01-19 16:25:59');
INSERT INTO `supplier_account` VALUES (3, 3, '工商银行上海分行', '6222000000000003', 'CNY', '91310100X00000003', 45, '2026-01-19 16:25:59', '2026-01-19 16:25:59');

-- ----------------------------
-- Table structure for supplier_contact
-- ----------------------------
DROP TABLE IF EXISTS `supplier_contact`;
CREATE TABLE `supplier_contact`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '联系人ID',
  `supplier_id` bigint UNSIGNED NOT NULL COMMENT '供应商ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '联系人姓名',
  `phone` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '电话',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '邮箱',
  `position` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '职位',
  `is_primary` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否主联系人',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_supplier_contact_supplier`(`supplier_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '供应商联系人' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of supplier_contact
-- ----------------------------
INSERT INTO `supplier_contact` VALUES (1, 1, '李经理', '13800001111', 'li@sup-a.com', '销售', 1, '2026-01-19 16:25:59', '2026-01-19 16:25:59');
INSERT INTO `supplier_contact` VALUES (2, 1, '王采购', '13800001112', 'wang@sup-a.com', '采购对接', 0, '2026-01-19 16:25:59', '2026-01-19 16:25:59');
INSERT INTO `supplier_contact` VALUES (3, 2, '周女士', '13800002221', 'zhou@pack.com', '商务', 1, '2026-01-19 16:25:59', '2026-01-19 16:25:59');
INSERT INTO `supplier_contact` VALUES (4, 3, '陈先生', '13800003331', 'chen@logi.com', '运营', 1, '2026-01-19 16:25:59', '2026-01-19 16:25:59');
INSERT INTO `supplier_contact` VALUES (5, 4, '赵小姐', '13800004441', 'zhao@elec.com', '销售', 1, '2026-01-19 16:25:59', '2026-01-19 16:25:59');

-- ----------------------------
-- Table structure for supplier_tag
-- ----------------------------
DROP TABLE IF EXISTS `supplier_tag`;
CREATE TABLE `supplier_tag`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '标签ID',
  `supplier_id` bigint UNSIGNED NOT NULL COMMENT '供应商ID',
  `tag` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '标签',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_supplier_tag`(`supplier_id` ASC, `tag` ASC) USING BTREE,
  INDEX `idx_supplier_tag_supplier`(`supplier_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '供应商标签' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of supplier_tag
-- ----------------------------
INSERT INTO `supplier_tag` VALUES (1, 1, '核心', '2026-01-19 16:25:59', '2026-01-19 16:25:59');
INSERT INTO `supplier_tag` VALUES (2, 1, '稳定', '2026-01-19 16:25:59', '2026-01-19 16:25:59');
INSERT INTO `supplier_tag` VALUES (3, 2, '包材', '2026-01-19 16:25:59', '2026-01-19 16:25:59');
INSERT INTO `supplier_tag` VALUES (4, 3, '物流', '2026-01-19 16:25:59', '2026-01-19 16:25:59');
INSERT INTO `supplier_tag` VALUES (5, 4, '电子', '2026-01-19 16:25:59', '2026-01-19 16:25:59');

-- ----------------------------
-- Table structure for supplier_type
-- ----------------------------
DROP TABLE IF EXISTS `supplier_type`;
CREATE TABLE `supplier_type`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '类型ID',
  `supplier_id` bigint UNSIGNED NOT NULL COMMENT '供应商ID',
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '类型(PRODUCT/PACKAGING/LOGISTICS)',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_supplier_type`(`supplier_id` ASC, `type` ASC) USING BTREE,
  INDEX `idx_supplier_type_supplier`(`supplier_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '供应商类型' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of supplier_type
-- ----------------------------
INSERT INTO `supplier_type` VALUES (1, 1, 'PRODUCT', '2026-01-19 16:25:59', '2026-01-19 16:25:59');
INSERT INTO `supplier_type` VALUES (2, 2, 'PACKAGING', '2026-01-19 16:25:59', '2026-01-19 16:25:59');
INSERT INTO `supplier_type` VALUES (3, 3, 'LOGISTICS', '2026-01-19 16:25:59', '2026-01-19 16:25:59');
INSERT INTO `supplier_type` VALUES (4, 4, 'PRODUCT', '2026-01-19 16:25:59', '2026-01-19 16:25:59');

-- ----------------------------
-- Table structure for system_log
-- ----------------------------
DROP TABLE IF EXISTS `system_log`;
CREATE TABLE `system_log`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `trace_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '追踪ID（全链路追踪，必填）',
  `level` enum('DEBUG','INFO','WARN','ERROR') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'INFO' COMMENT '日志级别',
  `module` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模块名（与审计日志保持一致）',
  `message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '日志消息',
  `context` json NULL COMMENT '上下文数据（请求参数、中间结果、变量状态等）',
  `exception` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '异常堆栈（ERROR级别时记录）',
  `file` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '触发日志的文件路径',
  `line` int UNSIGNED NULL DEFAULT NULL COMMENT '触发日志的行号',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_trace_id`(`trace_id` ASC) USING BTREE,
  INDEX `idx_level`(`level` ASC) USING BTREE,
  INDEX `idx_module`(`module` ASC) USING BTREE,
  INDEX `idx_gmt_create`(`gmt_create` ASC) USING BTREE,
  INDEX `idx_level_gmt_create`(`level` ASC, `gmt_create` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统日志表（错误、警告、调试信息）' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of system_log
-- ----------------------------

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名（登录名）',
  `password_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码哈希（bcrypt）',
  `real_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '真实姓名',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '手机号',
  `status` enum('ACTIVE','DISABLED') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ACTIVE' COMMENT '状态：ACTIVE=启用 DISABLED=禁用',
  `last_login_at` datetime NULL DEFAULT NULL COMMENT '最后登录时间',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_username`(`username` ASC) USING BTREE,
  INDEX `idx_email`(`email` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_gmt_create`(`gmt_create` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (8, 'admin', '$2a$10$UWjL9x/1Fc6/zslrOL.PyOdkEslSNAtMqKN36Tg4a84WaJG6x9x8u', '系统管理员', 'admin@example.com', NULL, 'ACTIVE', NULL, '2026-01-19 19:34:58', '2026-01-19 19:36:45');

-- ----------------------------
-- Table structure for user_role
-- ----------------------------
DROP TABLE IF EXISTS `user_role`;
CREATE TABLE `user_role`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `user_id` bigint UNSIGNED NOT NULL COMMENT '用户ID',
  `role_id` bigint UNSIGNED NOT NULL COMMENT '角色ID',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_user_id_role_id`(`user_id` ASC, `role_id` ASC) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_role_id`(`role_id` ASC) USING BTREE,
  INDEX `idx_gmt_create`(`gmt_create` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户角色关联表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_role
-- ----------------------------
INSERT INTO `user_role` VALUES (8, 8, 16, '2026-01-19 19:35:19', '2026-01-19 19:35:19');

-- ----------------------------
-- Table structure for warehouse
-- ----------------------------
DROP TABLE IF EXISTS `warehouse`;
CREATE TABLE `warehouse`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '仓库ID',
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '仓库代码（唯一标识）',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '仓库名称',
  `type` enum('FBA','THIRD_PARTY','OWN') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'OWN' COMMENT '仓库类型',
  `country` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '国家代码',
  `address` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '仓库地址',
  `contact_person` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '联系人',
  `contact_phone` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '联系电话',
  `contact_email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '联系邮箱',
  `status` enum('ACTIVE','INACTIVE','CLOSED') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ACTIVE' COMMENT '仓库状态',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '备注',
  `created_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint UNSIGNED NULL DEFAULT NULL COMMENT '更新人',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_code`(`code` ASC) USING BTREE,
  INDEX `idx_type`(`type` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_gmt_create`(`gmt_create` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '仓库表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of warehouse
-- ----------------------------
INSERT INTO `warehouse` VALUES (1, 'WH-CN-LOCAL', '本地仓库', 'OWN', 'CN', '', '', '', '', 'ACTIVE', '', 1, NULL, '2026-01-18 19:52:28', '2026-01-22 17:07:08');
INSERT INTO `warehouse` VALUES (2, 'WH-FBA-US', 'FBA  US 美国', 'FBA', 'US', '松仔园', '18410002000', '18410002000', '', 'ACTIVE', '', 1, NULL, '2026-01-18 19:52:28', '2026-01-28 14:11:16');
INSERT INTO `warehouse` VALUES (3, 'WH-FBA-AU', 'FBA  AU  澳洲', 'FBA', 'AU', '', '', '', '', 'ACTIVE', '', 1, NULL, '2026-01-18 19:52:28', '2026-01-22 17:05:17');

SET FOREIGN_KEY_CHECKS = 1;
