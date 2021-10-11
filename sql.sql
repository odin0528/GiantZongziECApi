ALTER TABLE `ec`.`platform` 
ADD COLUMN `fb_messenger_enabled` tinyint(4) NULL COMMENT 'FB Messenger 開關' AFTER `fb_page_id`;

ALTER TABLE `ec`.`platform` 
ADD COLUMN `icon_url` varchar(255) NULL AFTER `logo_url`;

CREATE TABLE `ec`.`platform_payment`  (
  `platform_id` int NOT NULL,
  `transfer_enabled` tinyint(255) NULL,
  `transfer_bank` varchar(255) NULL,
  `transfer_account` varchar(255) NULL,
  `delivery_enabled` tinyint(255) NULL,
  `delivery_711` tinyint(255) NULL,
  `delivery_family` tinyint(255) NULL,
  `delivery_hilife` tinyint(255) NULL,
  `delivery_ok` tinyint(255) NULL,
  `line_pay_enabled` tinyint(255) NULL,
  `updated_at` int NULL,
  `created_at` int NULL,
  PRIMARY KEY (`platform_id`)
);

CREATE DEFINER = `root`@`%` TRIGGER `insert_pages_rel` AFTER INSERT ON `platform` FOR EACH ROW BEGIN
	INSERT INTO pages (platform_id, url, title, is_menu, is_enabled, sort, created_at, updated_at, deleted_at, released_at) VALUES(NEW.id, 'index', '首頁', 1, 1, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0);
	INSERT INTO pages (platform_id, url, title, is_menu, is_enabled, sort, created_at, updated_at, deleted_at, released_at) VALUES(NEW.id, 'categories', '全部商品', 1, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0);
	INSERT INTO platform_payment (id, created_at, updated_at) VALUES(NEW.id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
END;

INSERT INTO `ec`.`platform_payment`(`platform_id`) VALUES (1);
INSERT INTO `ec`.`platform_payment`(`platform_id`) VALUES (2);
INSERT INTO `ec`.`platform_payment`(`platform_id`) VALUES (3);
INSERT INTO `ec`.`platform_payment`(`platform_id`) VALUES (4);
INSERT INTO `ec`.`platform_payment`(`platform_id`) VALUES (5);

ALTER TABLE `ec`.`platform` 
ADD COLUMN `description` text NULL COMMENT '網站簡介' AFTER `title`;