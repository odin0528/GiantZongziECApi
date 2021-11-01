ALTER TABLE `ec`.`platform_payment` 
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