ALTER TABLE `carts` 
ADD COLUMN `discount` double(10, 2) NULL COMMENT '折扣後的價格' AFTER `price`;

ALTER TABLE `order_products` 
ADD COLUMN `is_discount` tinyint(4) NULL AFTER `price`,
ADD COLUMN `discount` double(10, 2) NULL AFTER `is_discount`,
ADD COLUMN `discounted_price` double(10, 2) NULL AFTER `discount`;