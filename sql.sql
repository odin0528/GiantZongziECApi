ALTER TABLE `order_products` 
ADD COLUMN `is_discount` tinyint(4) NULL AFTER `price`,
ADD COLUMN `discount` double(10, 2) NULL AFTER `is_discount`,
ADD COLUMN `discounted_price` double(10, 2) NULL AFTER `discount`;

ALTER TABLE `ec`.`carts` 
DROP COLUMN `price`,
DROP COLUMN `total`,
DROP COLUMN `title`,
DROP COLUMN `style_title`,
DROP COLUMN `photo`,
DROP COLUMN `sku`;