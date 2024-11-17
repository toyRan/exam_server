/*
 Navicat Premium Dump SQL

 Source Server         : jx_ebook
 Source Server Type    : MySQL
 Source Server Version : 80012 (8.0.12)
 Source Host           : localhost:3306
 Source Schema         : jx_ebook

 Target Server Type    : MySQL
 Target Server Version : 80012 (8.0.12)
 File Encoding         : 65001

 Date: 02/11/2024 18:17:08
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for brands
-- ----------------------------
DROP TABLE IF EXISTS `brands`;
CREATE TABLE `brands`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `description` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of brands
-- ----------------------------

-- ----------------------------
-- Table structure for categories
-- ----------------------------
DROP TABLE IF EXISTS `categories`;
CREATE TABLE `categories`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `slug` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `pid` int(11) NOT NULL DEFAULT 0,
  `display_order` int(11) NOT NULL DEFAULT 50,
  `description` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of categories
-- ----------------------------
INSERT INTO `categories` VALUES (1, 'Optical', 'optical', 0, 50, 'dddd', '2024-10-29 17:15:16', '2024-10-31 16:06:16', NULL);
INSERT INTO `categories` VALUES (2, 'Sunglasses', 'sunglasses2', 0, 50, 'ddd', '2024-10-29 17:15:13', '2024-10-31 16:41:51', NULL);
INSERT INTO `categories` VALUES (3, '测试2222', 'test2222', 0, 0, '', '2024-10-29 17:14:57', '2024-10-29 17:15:03', '2024-10-29 17:16:35');
INSERT INTO `categories` VALUES (9, 'aa', 'dd', 1, 1, 'ddd', NULL, NULL, '2024-10-31 16:19:11');
INSERT INTO `categories` VALUES (10, 'test1', 'test1', 1, 1, '', NULL, NULL, NULL);

-- ----------------------------
-- Table structure for frame_materials
-- ----------------------------
DROP TABLE IF EXISTS `frame_materials`;
CREATE TABLE `frame_materials`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `description` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of frame_materials
-- ----------------------------
INSERT INTO `frame_materials` VALUES (1, 'TR99', '555', '2024-10-29 18:49:56', '2024-10-29 18:49:56', NULL);
INSERT INTO `frame_materials` VALUES (3, 'TR90', 'tr90', NULL, NULL, NULL);
INSERT INTO `frame_materials` VALUES (5, 'Acetate', 'Acetate', '2024-11-02 09:27:30', '2024-11-02 09:27:30', NULL);
INSERT INTO `frame_materials` VALUES (6, 'Metal', '', '2024-11-02 09:27:49', '2024-11-02 09:27:49', NULL);
INSERT INTO `frame_materials` VALUES (7, 'CP', '', '2024-11-02 09:27:57', '2024-11-02 09:27:57', NULL);
INSERT INTO `frame_materials` VALUES (8, 'PC', '', '2024-11-02 09:28:02', '2024-11-02 09:28:02', NULL);
INSERT INTO `frame_materials` VALUES (9, 'Titanium', '', '2024-11-02 09:28:13', '2024-11-02 09:28:13', NULL);
INSERT INTO `frame_materials` VALUES (10, 'Ultem', '', '2024-11-02 09:32:44', '2024-11-02 09:32:44', NULL);

-- ----------------------------
-- Table structure for password_reset_tokens
-- ----------------------------
DROP TABLE IF EXISTS `password_reset_tokens`;
CREATE TABLE `password_reset_tokens`  (
  `email` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `token` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`email`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of password_reset_tokens
-- ----------------------------

-- ----------------------------
-- Table structure for personal_access_tokens
-- ----------------------------
DROP TABLE IF EXISTS `personal_access_tokens`;
CREATE TABLE `personal_access_tokens`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `tokenable_type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `tokenable_id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `token` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `abilities` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `last_used_at` timestamp NULL DEFAULT NULL,
  `expires_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `personal_access_tokens_token_unique`(`token` ASC) USING BTREE,
  INDEX `personal_access_tokens_tokenable_type_tokenable_id_index`(`tokenable_type` ASC, `tokenable_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of personal_access_tokens
-- ----------------------------

-- ----------------------------
-- Table structure for product_images
-- ----------------------------
DROP TABLE IF EXISTS `product_images`;
CREATE TABLE `product_images`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `product_id` int(11) NULL DEFAULT NULL,
  `image_url` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `product_id`(`product_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of product_images
-- ----------------------------
INSERT INTO `product_images` VALUES (10, NULL, '1730278477107500300_延迟退休2.png');

-- ----------------------------
-- Table structure for products
-- ----------------------------
DROP TABLE IF EXISTS `products`;
CREATE TABLE `products`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `model_no` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `brand_id` int(11) NULL DEFAULT NULL,
  `series_id` int(11) NOT NULL,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `lens_width` int(11) NULL DEFAULT NULL,
  `nose_distance` int(11) NULL DEFAULT NULL,
  `temple_length` int(11) NULL DEFAULT NULL,
  `gender` enum('male','female','unisex') CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `description` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `category_id` int(11) NOT NULL,
  `sku_count` int(11) NULL DEFAULT 0 COMMENT 'sku数量',
  `status` tinyint(4) NOT NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `fk1`(`series_id` ASC) USING BTREE,
  INDEX `fk2`(`category_id` ASC) USING BTREE,
  INDEX `fk3`(`brand_id` ASC) USING BTREE,
  INDEX `deleted_at_index`(`deleted_at` ASC) USING BTREE,
  CONSTRAINT `fk1` FOREIGN KEY (`series_id`) REFERENCES `series` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `fk2` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `fk3` FOREIGN KEY (`brand_id`) REFERENCES `brands` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of products
-- ----------------------------

-- ----------------------------
-- Table structure for role_permission
-- ----------------------------
DROP TABLE IF EXISTS `role_permission`;
CREATE TABLE `role_permission`  (
  `role_id` int(11) NOT NULL,
  `permission_id` int(11) NOT NULL,
  PRIMARY KEY (`role_id`, `permission_id`) USING BTREE,
  INDEX `permission_id`(`permission_id` ASC) USING BTREE,
  CONSTRAINT `role_permission_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `role_permission_ibfk_2` FOREIGN KEY (`permission_id`) REFERENCES `sys_permissions` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of role_permission
-- ----------------------------

-- ----------------------------
-- Table structure for roles
-- ----------------------------
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `description` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of roles
-- ----------------------------
INSERT INTO `roles` VALUES (1, '超级管理员', NULL);
INSERT INTO `roles` VALUES (2, 'PDF添加员', NULL);

-- ----------------------------
-- Table structure for series
-- ----------------------------
DROP TABLE IF EXISTS `series`;
CREATE TABLE `series`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `description` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `pdf_url` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of series
-- ----------------------------
INSERT INTO `series` VALUES (1, 'MC101007', '', '', '2024-10-30 09:48:38', '2024-10-30 09:48:39', NULL);
INSERT INTO `series` VALUES (2, 'MC101010', '', '', '2024-10-30 09:48:51', '2024-10-30 09:48:52', NULL);
INSERT INTO `series` VALUES (3, 'MC101032', '', '', '2024-10-30 09:48:58', '2024-10-30 09:48:58', NULL);
INSERT INTO `series` VALUES (4, 'MC101033', '', '', '2024-10-30 09:49:03', '2024-10-30 09:49:03', NULL);
INSERT INTO `series` VALUES (5, 'MC21001', '', '', '2024-10-30 09:49:09', '2024-10-30 09:49:09', NULL);
INSERT INTO `series` VALUES (6, '333', '333', '', '2024-11-02 15:07:10', '2024-11-02 15:07:11', NULL);

-- ----------------------------
-- Table structure for sys_menus
-- ----------------------------
DROP TABLE IF EXISTS `sys_menus`;
CREATE TABLE `sys_menus`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `label` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '菜单名称',
  `parent_id` bigint(20) UNSIGNED NULL DEFAULT 0 COMMENT '父菜单ID',
  `order` int(11) NOT NULL DEFAULT 0 COMMENT '显示顺序',
  `link_to` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '' COMMENT '链接地址',
  `icon` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '#' COMMENT '菜单图标',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `menus_parent_id_foreign`(`parent_id` ASC) USING BTREE,
  CONSTRAINT `menus_parent_id_foreign` FOREIGN KEY (`parent_id`) REFERENCES `sys_menus` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 29 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_menus
-- ----------------------------
INSERT INTO `sys_menus` VALUES (1, '后台用户管理', NULL, 100, NULL, 'User', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (2, '后台用户列表', 1, 0, '/admin/sys-users', '#', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (3, '分类管理', NULL, 90, '', 'Operation', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (4, '分类列表', 3, 0, '/admin/categories', '#', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (5, '商品管理', NULL, 90, '', 'Goods', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (6, '商品列表', 5, 0, '/admin/products', '#', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (14, '后台角色管理', NULL, 99, '', 'UserFilled', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (15, '后台角色列表', 14, 0, '/admin/sys-roles', '#', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (16, '菜单管理', NULL, 95, '', 'Menu', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (17, '菜单列表', 16, 0, '/admin/menus', '#', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (18, '测试结果管理', NULL, 30, '', 'User', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (19, '测试结果列表', 18, 0, '/admin/exams', '#', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (20, '订单管理', NULL, 20, '', 'Coin', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (21, '订单列表', 20, 0, '/admin/orders', '#', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (22, '商品添加', 5, 0, '/admin/products/add', '#', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (23, '后台权限管理', NULL, 98, '', 'Coin', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (24, '后台权限列表', 23, 0, '/admin/sys-permissions', '', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (25, '框材质', 5, 0, '/admin/frame-materials', '#', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES (28, '系列Series', 5, 0, '/admin/series', '#', NULL, NULL, NULL);

-- ----------------------------
-- Table structure for sys_permissions
-- ----------------------------
DROP TABLE IF EXISTS `sys_permissions`;
CREATE TABLE `sys_permissions`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '权限名称',
  `description` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '权限描述',
  `route` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '匹配的路由 URL',
  `method` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '请求方法（GET, POST等）',
  `parent_id` int(11) NULL DEFAULT NULL COMMENT '父权限 ID（用于分层）',
  `created_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name` ASC) USING BTREE,
  INDEX `deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_permissions
-- ----------------------------
INSERT INTO `sys_permissions` VALUES (1, 'view_catalog', 'View product catalog', '222', 'GET', 0, '2024-11-01 16:55:56', '2024-11-01 17:30:13', NULL);
INSERT INTO `sys_permissions` VALUES (2, 'download_pdf', 'Download product PDF', '', 'GET', 1, NULL, NULL, NULL);
INSERT INTO `sys_permissions` VALUES (3, 'admin_manage_users', 'Manage users', '', 'POST', NULL, NULL, NULL, NULL);
INSERT INTO `sys_permissions` VALUES (4, 'admin_manage_catalog', 'Manage catalog', '', 'POST', NULL, NULL, NULL, NULL);
INSERT INTO `sys_permissions` VALUES (5, '查看最新款', 'xxx', '1', 'GET', 0, '2024-11-01 17:34:58', '2024-11-01 17:34:58', NULL);
INSERT INTO `sys_permissions` VALUES (6, '测试1', '测试111', '1', 'GET', 0, '2024-11-01 17:34:46', '2024-11-01 17:34:46', NULL);

-- ----------------------------
-- Table structure for sys_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_menu`;
CREATE TABLE `sys_role_menu`  (
  `sys_role_id` bigint(20) NOT NULL COMMENT '角色ID',
  `sys_menu_id` bigint(20) NOT NULL COMMENT '菜单ID',
  PRIMARY KEY (`sys_role_id`, `sys_menu_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '角色和菜单关联表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role_menu
-- ----------------------------
INSERT INTO `sys_role_menu` VALUES (1, 1);
INSERT INTO `sys_role_menu` VALUES (1, 2);
INSERT INTO `sys_role_menu` VALUES (1, 3);
INSERT INTO `sys_role_menu` VALUES (1, 4);
INSERT INTO `sys_role_menu` VALUES (1, 5);
INSERT INTO `sys_role_menu` VALUES (1, 6);
INSERT INTO `sys_role_menu` VALUES (1, 7);
INSERT INTO `sys_role_menu` VALUES (1, 8);
INSERT INTO `sys_role_menu` VALUES (1, 9);
INSERT INTO `sys_role_menu` VALUES (1, 10);
INSERT INTO `sys_role_menu` VALUES (1, 14);
INSERT INTO `sys_role_menu` VALUES (1, 15);
INSERT INTO `sys_role_menu` VALUES (1, 23);
INSERT INTO `sys_role_menu` VALUES (1, 24);
INSERT INTO `sys_role_menu` VALUES (1, 25);
INSERT INTO `sys_role_menu` VALUES (1, 28);

-- ----------------------------
-- Table structure for sys_role_permission
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_permission`;
CREATE TABLE `sys_role_permission`  (
  `sys_role_id` int(11) NOT NULL,
  `sys_permission_id` int(11) NOT NULL,
  PRIMARY KEY (`sys_role_id`, `sys_permission_id`) USING BTREE,
  INDEX `fk2`(`sys_permission_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role_permission
-- ----------------------------
INSERT INTO `sys_role_permission` VALUES (1, 1);
INSERT INTO `sys_role_permission` VALUES (1, 2);
INSERT INTO `sys_role_permission` VALUES (1, 3);
INSERT INTO `sys_role_permission` VALUES (1, 4);

-- ----------------------------
-- Table structure for sys_roles
-- ----------------------------
DROP TABLE IF EXISTS `sys_roles`;
CREATE TABLE `sys_roles`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `description` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name` ASC) USING BTREE,
  INDEX `deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_roles
-- ----------------------------
INSERT INTO `sys_roles` VALUES (1, '超级管理员', 'cxxx', '2024-10-31 17:22:05', '2024-10-31 17:22:23', NULL);
INSERT INTO `sys_roles` VALUES (3, 'PDF添加员', 'xxx', '2024-10-31 17:22:10', '2024-10-31 17:22:13', NULL);
INSERT INTO `sys_roles` VALUES (4, '测试员', '111222', '2024-10-31 17:16:07', '2024-10-31 17:19:55', '2024-10-31 17:48:52');
INSERT INTO `sys_roles` VALUES (5, '测试员2', 'dddd2222222', '2024-10-31 17:16:42', '2024-10-31 17:19:52', '2024-10-31 17:48:48');

-- ----------------------------
-- Table structure for sys_user_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_user_role`;
CREATE TABLE `sys_user_role`  (
  `sys_user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `sys_role_id` bigint(20) NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`sys_user_id`, `sys_role_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '用户和角色关联表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_user_role
-- ----------------------------
INSERT INTO `sys_user_role` VALUES (1, 1);
INSERT INTO `sys_user_role` VALUES (3, 1);
INSERT INTO `sys_user_role` VALUES (4, 3);
INSERT INTO `sys_user_role` VALUES (5, 1);
INSERT INTO `sys_user_role` VALUES (6, 1);
INSERT INTO `sys_user_role` VALUES (7, 1);
INSERT INTO `sys_user_role` VALUES (7, 3);
INSERT INTO `sys_user_role` VALUES (8, 1);
INSERT INTO `sys_user_role` VALUES (8, 3);
INSERT INTO `sys_user_role` VALUES (9, 1);
INSERT INTO `sys_user_role` VALUES (9, 3);
INSERT INTO `sys_user_role` VALUES (11, 1);
INSERT INTO `sys_user_role` VALUES (13, 1);
INSERT INTO `sys_user_role` VALUES (13, 3);
INSERT INTO `sys_user_role` VALUES (14, 1);
INSERT INTO `sys_user_role` VALUES (14, 3);

-- ----------------------------
-- Table structure for sys_users
-- ----------------------------
DROP TABLE IF EXISTS `sys_users`;
CREATE TABLE `sys_users`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `email` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `password` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `username`(`username` ASC) USING BTREE,
  UNIQUE INDEX `email`(`email` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 15 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_users
-- ----------------------------
INSERT INTO `sys_users` VALUES (3, 'abc3', 'abc3@qq.com', '$2a$10$eaZCe9mwlrfOvtu.M2MJcumiyZr/HgxvAQYdAVI/cSFOG8CNglstC', 0, '2024-10-31 11:57:19', '2024-10-31 15:16:59', NULL);
INSERT INTO `sys_users` VALUES (5, 'q1', 'q1@qq.com', '$2a$10$Rt5COWFlJ2BBjzT2.a0at.o9JXo0Ty0wqOgco1iYYIx50VkMbfSU2', 0, '2024-10-31 13:35:48', '2024-10-31 14:01:27', '2024-10-31 16:03:00');
INSERT INTO `sys_users` VALUES (6, '张三', 'zhangsan@qq.com', '$2a$10$NTm.2d463hJJcrgG77bL2.WgM4Xez8T0gVxKXQJSyPaWNxu18qTE2', 0, '2024-10-31 14:02:55', '2024-10-31 14:03:22', '2024-10-31 16:03:04');
INSERT INTO `sys_users` VALUES (7, 'ffff', 'fff@qq.com', '$2a$10$1ZEp0DrMJuEa6tythGB4WOWtMgjpzD5GsP63qNdlDnRQZskVnOtT.', 1, '2024-10-31 14:06:09', '2024-10-31 16:03:40', NULL);
INSERT INTO `sys_users` VALUES (8, 'ddd', 'ddd@qq.com', '$2a$10$QsFw/xH6SqzYscCSzh5TeeZjIQtsGUQ/l3Msi/M9yS8XVRC4zdlQ2', 0, '2024-10-31 14:12:15', '2024-10-31 14:12:51', NULL);
INSERT INTO `sys_users` VALUES (9, 'dddd3', 'dddd@qq.com', '$2a$10$ULiJVKL1B2WpEnfc/aCCbuihRdRN1E8BOoxvasCwSRsV3vlNogC1G', 1, '2024-10-31 14:21:49', '2024-10-31 14:21:50', NULL);
INSERT INTO `sys_users` VALUES (11, 'sss2', 'sss2@qq.com', '$2a$10$fVSN5COyJAaw4H.5CQkgk.KW1uTMq2C8G/KWh/MGR1KVK3oHY/V7K', 1, '2024-10-31 14:25:51', '2024-10-31 14:41:32', NULL);
INSERT INTO `sys_users` VALUES (13, 'diedie', 'diedie@qq.com', '$2a$10$1t2NYOueLefsArxfqwICGeRZ8PHeXKTJYBR9YAVYKYUgUPPxwfZh6', 0, '2024-10-31 15:17:46', '2024-10-31 15:17:47', NULL);
INSERT INTO `sys_users` VALUES (14, 'hehe', 'laoliu@qq.com', '$2a$10$6g5MDalso7UO3bhkEw9Yoer.IUoQjTyXZopYZ/qAH.hoEezJcFZaK', 0, '2024-10-31 15:21:10', '2024-10-31 15:21:23', '2024-10-31 16:03:06');

-- ----------------------------
-- Table structure for user_role
-- ----------------------------
DROP TABLE IF EXISTS `user_role`;
CREATE TABLE `user_role`  (
  `user_id` int(11) NOT NULL,
  `role_id` int(11) NOT NULL,
  PRIMARY KEY (`user_id`, `role_id`) USING BTREE,
  INDEX `role_id`(`role_id` ASC) USING BTREE,
  CONSTRAINT `user_role_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `user_role_ibfk_2` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_role
-- ----------------------------
INSERT INTO `user_role` VALUES (1, 1);

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `first_name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `last_name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `email` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `password` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `email`(`email` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (1, 'gao', 'sheng', 'abc@qq.com', '$2a$10$I4gYU.w3sAHCpPcHTAuur.t67LPK4lcWWBLKu/qqqqbVeyS.w/x42', 1, '2024-10-28 16:11:30', '2024-10-28 16:11:30', NULL);
INSERT INTO `users` VALUES (2, 'gao', 'sheng', 'abc2@qq.com', '$2a$10$rbcXRMXQjY5B.GcgtaclReIiJ.PZHlSpwygq48AG3kbQZIaKyCsFu', 1, '2024-10-28 16:20:29', '2024-10-28 16:20:29', NULL);

SET FOREIGN_KEY_CHECKS = 1;
