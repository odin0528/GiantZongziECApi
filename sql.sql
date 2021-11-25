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
