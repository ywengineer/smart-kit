/*
 Navicat Premium Dump SQL

 Source Server         : local-mysql
 Source Server Type    : MySQL
 Source Server Version : 80040 (8.0.40)
 Source Host           : localhost:3306
 Source Schema         : smart-passport

 Target Server Type    : MySQL
 Target Server Version : 80040 (8.0.40)
 File Encoding         : 65001

 Date: 06/02/2025 13:46:06
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for log_purchase
-- ----------------------------
DROP TABLE IF EXISTS `log_purchase`;
CREATE TABLE `log_purchase`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `channel` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '支付渠道',
  `quantity` tinyint NULL DEFAULT NULL,
  `transaction_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `original_transaction_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `purchase_date` timestamp NULL DEFAULT NULL COMMENT '充值时间',
  `original_purchase_date` timestamp NULL DEFAULT NULL COMMENT '第一次订阅时间',
  `expired_date` timestamp NULL DEFAULT NULL COMMENT '服务过期时间',
  `app_item_id` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT '标识程序的字符串或purchaseToken',
  `original_application_version` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT 'AppVersion或者developerPayload',
  `version_external_identifier` int NULL DEFAULT NULL,
  `price` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '价格',
  `free_trail` bit(1) NOT NULL DEFAULT b'0' COMMENT '试用订单',
  `test_order` bit(1) NOT NULL DEFAULT b'0' COMMENT '测试订单',
  `receipt` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT '订单原始凭证',
  `receiptResult` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT '订单验证结果',
  `bundle_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `game_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '游戏ID',
  `server_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '服务器ID',
  `passport` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '账号',
  `player_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '玩家ID',
  `player_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '玩家名称',
  `rankPoints` int UNSIGNED NULL DEFAULT 0 COMMENT '段位积分',
  `credits` bigint UNSIGNED NULL DEFAULT 0 COMMENT '点券剩余量',
  `money` bigint UNSIGNED NULL DEFAULT 0 COMMENT '钻石剩余量',
  `coin` bigint UNSIGNED NULL DEFAULT 0 COMMENT '当前金币',
  `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '日志时间',
  `notifyTimes` int NOT NULL DEFAULT 1 COMMENT '通知次数',
  `notified` bit(1) NOT NULL DEFAULT b'1' COMMENT '是否通知游戏服',
  `finishTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '订单处理结束时间',
  `locale` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '地区',
  `systemType` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '系统类型',
  `playerCreateTime` timestamp NULL DEFAULT NULL COMMENT '角色创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `ix_order_id`(`transaction_id` ASC) USING BTREE,
  INDEX `index_notified`(`notified` ASC) USING BTREE,
  INDEX `index_passport`(`passport` ASC, `player_id` ASC, `player_name` ASC) USING BTREE,
  INDEX `index_product_player`(`product_id` ASC, `server_id` ASC, `passport` ASC, `player_id` ASC, `player_name` ASC) USING BTREE,
  INDEX `index_server_player`(`server_id` ASC, `passport` ASC, `player_id` ASC, `player_name` ASC) USING BTREE,
  INDEX `index_transaction`(`original_transaction_id` ASC, `transaction_id` ASC) USING BTREE,
  INDEX `ix_player_time`(`player_id` ASC, `purchase_date` ASC) USING BTREE,
  INDEX `ix_game_server_prod`(`game_id` ASC, `server_id` ASC, `product_id` ASC) USING BTREE,
  INDEX `ix_time`(`purchase_date` ASC) USING BTREE,
  INDEX `ix_game_server_player`(`game_id` ASC, `server_id` ASC, `player_id` ASC, `test_order` ASC, `free_trail` ASC) USING BTREE,
  INDEX `ix_game_server_purchase`(`game_id` ASC, `server_id` ASC, `purchase_date` ASC, `test_order` ASC, `free_trail` ASC) USING BTREE,
  INDEX `ix_game_server_player_time`(`game_id` ASC, `player_id` ASC, `server_id` ASC, `purchase_date` ASC, `test_order` ASC, `free_trail` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for passport_bindings
-- ----------------------------
DROP TABLE IF EXISTS `passport_bindings`;
CREATE TABLE `passport_bindings`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `passport_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '账号ID',
  `bind_type` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '绑定类型',
  `bind_id` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '绑定第3方平台ID',
  `access_token` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `refresh_token` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `social_name` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `gender` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `icon_url` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_passport_bindings_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_passport_id`(`passport_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for passport_punishes
-- ----------------------------
DROP TABLE IF EXISTS `passport_punishes`;
CREATE TABLE `passport_punishes`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `passport_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '账号ID',
  `device_id` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '设备ID',
  `type` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '惩罚类型',
  `begin_time` datetime NULL DEFAULT NULL,
  `end_time` datetime NULL DEFAULT NULL,
  `reason` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '原因',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_passport_punishes_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_passport`(`passport_id` ASC) USING BTREE,
  INDEX `idx_device`(`device_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for passports
-- ----------------------------
DROP TABLE IF EXISTS `passports`;
CREATE TABLE `passports`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `device_id` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '设备ID',
  `adid` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '设备广告标识',
  `system_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '系统类型',
  `locale` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '地区',
  `extra` json NULL COMMENT '主机额外信息',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_passports_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_device`(`device_id` ASC) USING BTREE,
  INDEX `idx_adid`(`adid` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for subscription
-- ----------------------------
DROP TABLE IF EXISTS `subscription`;
CREATE TABLE `subscription`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `game_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '游戏ID',
  `server_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '服务器ID',
  `passport` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '账号',
  `player_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '玩家ID',
  `player_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '玩家名称',
  `original_transaction_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT 'CURRENT_TIMESTAMP' COMMENT '订阅ID',
  `original_purchase_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '第一次订阅时间',
  `product_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '产品ID',
  `channel` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT 'GooglePlay' COMMENT '支付渠道',
  `bundle_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT 'appid',
  `receipt` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '支付原始凭据',
  `receiptResult` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '凭据明文数据',
  `auto_renew_status` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '1' COMMENT '是否自动续期',
  `expires_date` bigint NOT NULL DEFAULT 0 COMMENT '过期时间',
  `expiration_intent` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '过期原因',
  `is_in_billing_retry_period` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '是否在重试自动订阅',
  `cancellation_reason` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '取消订阅原因',
  `linkedPurchaseToken` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT 'Google订阅原始凭据',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uix_game_server_player_prod`(`game_id` ASC, `server_id` ASC, `player_id` ASC, `product_id` ASC) USING BTREE,
  INDEX `uix_subs`(`game_id` ASC, `server_id` ASC, `player_id` ASC, `product_id` ASC, `channel` ASC) USING BTREE,
  INDEX `ix_transaction`(`original_transaction_id`(191) ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '商品订阅数据' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for white_lists
-- ----------------------------
DROP TABLE IF EXISTS `white_lists`;
CREATE TABLE `white_lists`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `passport` bigint UNSIGNED NULL DEFAULT NULL COMMENT '账号ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uni_white_lists_passport`(`passport` ASC) USING BTREE,
  INDEX `idx_white_lists_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
