ALTER TABLE `platform_payment` 
DROP COLUMN `transfer_bank`,
DROP COLUMN `transfer_account`,
CHANGE COLUMN `transfer_enabled` `webatm_enabled` tinyint(4) NULL DEFAULT NULL AFTER `line_pay_enabled`,
MODIFY COLUMN `delivery_enabled` tinyint(4) NULL DEFAULT NULL AFTER `platform_id`,
MODIFY COLUMN `delivery_711` tinyint(4) NULL DEFAULT NULL AFTER `delivery_enabled`,
MODIFY COLUMN `delivery_family` tinyint(4) NULL DEFAULT NULL AFTER `delivery_711`,
MODIFY COLUMN `delivery_hilife` tinyint(4) NULL DEFAULT NULL AFTER `delivery_family`,
MODIFY COLUMN `delivery_ok` tinyint(4) NULL DEFAULT NULL AFTER `delivery_hilife`,
ADD COLUMN `credit_card_enabled` tinyint(4) NULL AFTER `delivery_ok`,
MODIFY COLUMN `line_pay_enabled` tinyint(4) NULL DEFAULT NULL AFTER `delivery_ok`,
ADD COLUMN `atm_enabled` tinyint(4) NULL AFTER `webatm_enabled`,
ADD COLUMN `cvs_enabled` tinyint(4) NULL AFTER `atm_enabled`,
ADD COLUMN `barcode_enabled` tinyint(4) NULL AFTER `cvs_enabled`;

ALTER TABLE `orders` 
ADD COLUMN `ecpay_mac` varchar(255) NULL COMMENT '綠界檢查碼' AFTER `transaction_id`;

ALTER TABLE `orders` 
DROP COLUMN `ecpay_mac`;

ALTER TABLE `orders` 
ADD COLUMN `shipment_no` varchar(63) NULL COMMENT '託運單號' AFTER `store_phone`,
ADD COLUMN `logistics_status` varchar(255) NULL COMMENT '物流狀態' AFTER `shipment_no`,
ADD COLUMN `logistics_msg` varchar(255) NULL COMMENT '物流狀態說明' AFTER `logistics_status`;

ALTER TABLE `orders` 
MODIFY COLUMN `logistics_status` tinyint(4) NULL DEFAULT NULL COMMENT '物流狀態' AFTER `shipment_no`;

ALTER TABLE `orders` 
ADD COLUMN `logistics_id` varchar(63) NULL COMMENT '物流單號' AFTER `store_phone`
MODIFY COLUMN `logistics_id` varchar(63) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '物流單號' AFTER `store_phone`,
MODIFY COLUMN `shipment_no` varchar(63) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '託運單號' AFTER `logistics_id`,
MODIFY COLUMN `logistics_status` tinyint(4) NULL DEFAULT 0 COMMENT '物流狀態' AFTER `shipment_no`,
MODIFY COLUMN `logistics_msg` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '物流狀態說明' AFTER `logistics_status`;

ALTER TABLE `orders` 
ADD COLUMN `county` varchar(255) NULL COMMENT '宅配地址(縣)' AFTER `fullname`,
ADD COLUMN `district` varchar(255) NULL COMMENT '宅配地址(地區)' AFTER `county`,
ADD COLUMN `zip_code` varchar(15) NULL COMMENT '宅配郵遞區號' AFTER `district`;

ALTER TABLE `member_delivery` 
ADD COLUMN `county` varchar(255) NULL COMMENT '宅配地址(縣)' AFTER `fullname`,
ADD COLUMN `district` varchar(255) NULL COMMENT '宅配地址(地區)' AFTER `county`,
ADD COLUMN `zip_code` varchar(15) NULL COMMENT '宅配郵遞區號' AFTER `district`;

ALTER TABLE `orders` 
ADD COLUMN `payment_charge_fee` double(10, 2) NULL COMMENT '金流交易手續費' AFTER `qty`,
ADD COLUMN `logistics_charge_fee` double(10, 2) NULL COMMENT '實際物流運費' AFTER `payment_charge_fee`,
MODIFY COLUMN `payment` tinyint(255) NULL DEFAULT NULL COMMENT '付費方式(2:貨到付款 3:信用卡 4:line pay 5:atm 6:超商代碼 7:超商條碼 )' AFTER `memo`

CREATE TABLE `platform_logistics`  (
  `platform_id` int NOT NULL,
  `home_enabled` tinyint(4) NULL,
  `home_charge_fee` int(255) NULL,
  `uni_enabled` tinyint(4) NULL,
  `uni_charge_fee` int(255) NULL,
  `family_enabled` tinyint(4) NULL,
  `family_charge_fee` int(255) NULL,
  `hilife_enabled` tinyint(4) NULL,
  `hilife_charge_fee` int(255) NULL,
  `ok_enabled` tinyint(4) NULL,
  `ok_charge_fee` int(255) NULL,
  `created_at` int(255) NULL,
  `updated_at` int(255) NULL,
  PRIMARY KEY (`platform_id`)
);

DROP TRIGGER `insert_pages_rel`;
CREATE DEFINER = `root`@`%` TRIGGER `insert_pages_rel` AFTER INSERT ON `platform` FOR EACH ROW BEGIN
	INSERT INTO pages (platform_id, url, title, is_enabled, created_at, updated_at, deleted_at, released_at) VALUES(NEW.id, 'index', '首頁', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0);
	INSERT INTO pages (platform_id, url, title, is_enabled, created_at, updated_at, deleted_at, released_at) VALUES(NEW.id, 'categories', '全部商品', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0);
	INSERT INTO platform_menu (platform_id, title, link_type, link, sort, is_enabled, created_at, updated_at, deleted_at) VALUES(NEW.id, '首頁', 1, 'index', 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0);
	INSERT INTO platform_menu (platform_id, title, link_type, link, sort, is_enabled, created_at, updated_at, deleted_at) VALUES(NEW.id, '全部商品', 1, 'categories', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0);
	INSERT INTO platform_payment (platform_id, created_at, updated_at) VALUES(NEW.id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
	INSERT INTO platform_logistics (platform_id, created_at, updated_at) VALUES(NEW.id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
END;

ALTER TABLE `platform_payment` 
DROP COLUMN `delivery_711`,
DROP COLUMN `delivery_family`,
DROP COLUMN `delivery_hilife`,
DROP COLUMN `delivery_ok`,
DROP COLUMN `webatm_enabled`;

insert into platform_logistics (platform_id, uni_enabled, uni_charge_fee) (
	SELECT id, 1, 60 FROM platform
);