-- Adminer 4.8.1 MySQL 8.0.26 dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

CREATE DATABASE `ics_test_db` /*!40100 DEFAULT CHARACTER SET utf8 */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `ics_test_db`;

DROP TABLE IF EXISTS `bank`;
CREATE TABLE `bank` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(200) NOT NULL,
  `account_no` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

INSERT INTO `bank` (`id`, `name`, `account_no`) VALUES
(1,	'KTB',	'999-999-9999');

DROP TABLE IF EXISTS `category`;
CREATE TABLE `category` (
  `id` int NOT NULL AUTO_INCREMENT,
  `category` varchar(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

INSERT INTO `category` (`id`, `category`) VALUES
(1,	'Plain Color'),
(2,	'Pattern'),
(3,	'Figure');

DROP TABLE IF EXISTS `order_list`;
CREATE TABLE `order_list` (
  `order_id` int NOT NULL,
  `product_id` int NOT NULL,
  `amount` int NOT NULL,
  `price_per_unit` float NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

INSERT INTO `order_list` (`order_id`, `product_id`, `amount`, `price_per_unit`) VALUES
(31,	9,	2,	100),
(31,	11,	1,	175);

DROP TABLE IF EXISTS `order_status`;
CREATE TABLE `order_status` (
  `id` int NOT NULL AUTO_INCREMENT,
  `status` varchar(30) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

INSERT INTO `order_status` (`id`, `status`) VALUES
(1,	'Pending Payment'),
(2,	'Delivery'),
(3,	'Succeed'),
(4,	'Cancel');

DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(30) NOT NULL,
  `total_price` float NOT NULL DEFAULT '0',
  `address` varchar(255) NOT NULL,
  `status_id` tinyint NOT NULL,
  `payment_id` int NOT NULL,
  `create_at` int NOT NULL DEFAULT '0',
  `update_at` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

INSERT INTO `orders` (`id`, `user_id`, `total_price`, `address`, `status_id`, `payment_id`, `create_at`, `update_at`) VALUES
(31,	'dXNlcjoxMjM0NTY=',	375,	'123 จ.ลำปาง',	3,	7,	1642067817,	1642067969);

DROP TABLE IF EXISTS `payment`;
CREATE TABLE `payment` (
  `id` int NOT NULL AUTO_INCREMENT,
  `bank_id` int NOT NULL,
  `price` float NOT NULL,
  `receipt_path` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `create_at` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

INSERT INTO `payment` (`id`, `bank_id`, `price`, `receipt_path`, `create_at`) VALUES
(7,	1,	375,	'/storage/receipt/2022/01/13/1642067896.jpeg',	1642067919);

DROP TABLE IF EXISTS `product`;
CREATE TABLE `product` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID สินค้า',
  `name` varchar(255) NOT NULL COMMENT 'ชื่อสินค้า',
  `gender` varchar(6) NOT NULL DEFAULT '0' COMMENT 'male, female',
  `category_id` tinyint NOT NULL DEFAULT '0',
  `size_id` int NOT NULL DEFAULT '0',
  `total_amount` int NOT NULL DEFAULT '0',
  `price_per_unit` float NOT NULL DEFAULT '0',
  `create_at` int NOT NULL DEFAULT '0',
  `update_at` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

INSERT INTO `product` (`id`, `name`, `gender`, `category_id`, `size_id`, `total_amount`, `price_per_unit`, `create_at`, `update_at`) VALUES
(9,	'เสื้อฮาวาย',	'male',	1,	1,	8,	100,	1642067485,	1642067817),
(10,	'เสื้อฮาวาย',	'male',	1,	3,	10,	150,	1642067508,	1642067508),
(11,	'ชุดเดรท',	'female',	1,	2,	14,	175,	1642067544,	1642067817),
(12,	'ชุดเดรท',	'female',	1,	1,	15,	175,	1642067556,	1642067556);

DROP TABLE IF EXISTS `size`;
CREATE TABLE `size` (
  `id` int NOT NULL AUTO_INCREMENT,
  `size` varchar(2) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

INSERT INTO `size` (`id`, `size`) VALUES
(1,	'XS'),
(2,	'S'),
(3,	'M'),
(4,	'L'),
(5,	'XL');

-- 2022-01-13 10:02:00
