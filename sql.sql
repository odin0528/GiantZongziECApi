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