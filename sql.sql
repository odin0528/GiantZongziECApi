ALTER TABLE `order_products` 
ADD COLUMN `is_discount` tinyint(4) NULL AFTER `price`,
ADD COLUMN `discount` double(10, 2) NULL AFTER `is_discount`,
ADD COLUMN `discounted_price` double(10, 2) NULL AFTER `discount`;

ALTER TABLE `carts` 
DROP COLUMN `price`,
DROP COLUMN `total`,
DROP COLUMN `title`,
DROP COLUMN `style_title`,
DROP COLUMN `photo`,
DROP COLUMN `sku`;

ALTER TABLE `product_style_table` 
CHANGE COLUMN `sub_title` `sub_style_title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL AFTER `title`,
ADD COLUMN `style_title` varchar(255) NULL AFTER `title`;

ALTER TABLE `product_style_table` 
ADD COLUMN `photo` varchar(255) NULL AFTER `sub_style_title`;

update product_style_table set style_title = title;
update product_style_table set title = (select title from products where id = product_style_table.product_id);
UPDATE product_style_table set photo = (select img from product_style where product_style.id = (select id from product_style where product_id = product_style_table.product_id and product_style.title = product_style_table.style_title))
UPDATE product_style_table set photo = (select img from product_photos where product_photos.product_id = product_style_table.product_id and sort = 1) where photo = '' or photo is null;


ALTER TABLE `products` 
ADD COLUMN `min` double(10, 2) NULL COMMENT '所有規格中最低的價格' AFTER `sub_style_enabled`,
ADD COLUMN `max` double(10, 2) NULL COMMENT '所有規格中最高的價格' AFTER `min`,
ADD COLUMN `photo` varchar(255) NULL COMMENT '首圖' AFTER `max`;

UPDATE products SET 
min = (SELECT price FROM product_style_table WHERE product_style_table.product_id = products.id ORDER BY price ASC LIMIT 1),
max = (SELECT price FROM product_style_table WHERE product_style_table.product_id = products.id ORDER BY price DESC LIMIT 1);

ALTER TABLE `product_style_table` 
CHANGE COLUMN `group` `group_no` tinyint(255) NULL DEFAULT NULL AFTER `product_id`;

ALTER TABLE `orders` 
ADD COLUMN `paid_at` int(11) NULL COMMENT '付款時間' AFTER `status`,
ADD COLUMN `delivered_at` int(11) NULL COMMENT '到貨時間' AFTER `paid_at`,
ADD COLUMN `picker_id` int(11) NULL DEFAULT 0 COMMENT '撿貨者id' AFTER `member_id`
MODIFY COLUMN `created_at` int(11) NULL DEFAULT NULL AFTER `status`;

UPDATE orders SET `status` = 51 WHERE `status` = 31;
UPDATE orders SET `status` = 99 WHERE `status` = -99;
UPDATE orders SET `logistics_status` = 110 WHERE `logsitics_status` = 3;
UPDATE orders SET `logistics_status` = 120 WHERE `logsitics_status` = 4;
UPDATE orders SET `logistics_status` = 199 WHERE `logsitics_status` = 5;