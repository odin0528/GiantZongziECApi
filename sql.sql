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
ADD COLUMN `picker_id` int(11) NULL DEFAULT 0 COMMENT '撿貨者id' AFTER `member_id`,
MODIFY COLUMN `logistics_status` int(11) UNSIGNED NULL DEFAULT 0 COMMENT '物流狀態' AFTER `shipment_no`,
MODIFY COLUMN `created_at` int(11) NULL DEFAULT NULL AFTER `status`;

UPDATE orders SET `status` = 51 WHERE `status` = 31;
UPDATE orders SET `status` = 51 WHERE `status` = 31;
UPDATE orders SET `status` = 99 WHERE `status` = -99;
UPDATE orders SET `logistics_status` = 110 WHERE `logistics_status` = 3;
UPDATE orders SET `logistics_status` = 120 WHERE `logistics_status` = 4;
UPDATE orders SET `logistics_status` = 199 WHERE `logistics_status` = 5;
UPDATE orders SET `status` = 24 WHERE `status` = 21 and shipment_no is not null;

ALTER TABLE `platform` 
ADD COLUMN `mobile_logo_url` varchar(255) NULL AFTER `logo_url`;

ALTER TABLE `product_style_table` 
MODIFY COLUMN `group_no` tinyint(4) UNSIGNED NULL DEFAULT NULL AFTER `product_id`,
ADD COLUMN `sort` tinyint(4) UNSIGNED NULL AFTER `group_no`;

ALTER TABLE `product_style_table` 
ADD COLUMN `cost` double(10, 2) NULL COMMENT '單件產品成本' AFTER `qty`,
ADD COLUMN `suggest_price` double(10, 2) NULL COMMENT '建議售價' AFTER `cost`,
ADD COLUMN `no_store_delivery` tinyint(4) NULL COMMENT '超過此數量後不可超取，若為0則皆可超取' AFTER `suggest_price`,
ADD COLUMN `no_over_sale` tinyint(4) NULL COMMENT '不可超賣' AFTER `no_store_delivery`;

ALTER TABLE `carts` 
CHANGE COLUMN `qty` `buy_count` int(11) NOT NULL AFTER `style_id`;

ALTER TABLE `products` 
ADD COLUMN `sold` int(11) UNSIGNED NULL COMMENT '賣出數量' AFTER `photo`;


ALTER TABLE `product_style_table` 
ADD COLUMN `sold` int(11) UNSIGNED NULL COMMENT '賣出數量' AFTER `no_over_sale`;

CREATE TRIGGER `sold_out` BEFORE INSERT ON `order_products` FOR EACH ROW BEGIN
	DECLARE affected tinyint;
	DECLARE msg varchar(128);
	
	UPDATE product_style_table 
	SET qty = qty - new.qty, sold = sold + new.qty
	WHERE id = new.style_id AND (no_over_sale = 0 OR qty >= new.qty);
	SELECT ROW_COUNT() into affected;
	if affected = 0 then
		set msg = 'out_of_stock';
    signal sqlstate '45000' set message_text = msg;
	end if;
	
	UPDATE products SET sold = sold + new.qty WHERE new.product_id;
END;

UPDATE product_style_table SET sold = (select count(*) from order_products where order_products.style_id = product_style_table.id);
UPDATE products SET sold = (select count(*) from order_products where order_products.product_id = products.id);
UPDATE product_style_table SET no_over_sale = 0 WHERE no_over_sale IS NULL;
UPDATE product_style_table SET no_store_delivery = 0 WHERE no_store_delivery IS NULL;
UPDATE product_style_table SET cost = 0 WHERE cost IS NULL;
UPDATE product_style_table SET suggest_price = 0 WHERE suggest_price IS NULL;


ALTER TABLE `product_style_table` 
ADD COLUMN `low_stock` int(11) NULL COMMENT '低庫存數量，庫存低於此數字時發通知' AFTER `qty`;

UPDATE product_style_table SET low_stock = 0;

CREATE VIEW `report_low_stock` AS SELECT
	`products`.`id` AS `id`,
	`products`.`platform_id` AS `platform_id`,
	`products`.`category_layer1` AS `category_layer1`,
	`products`.`category_layer2` AS `category_layer2`,
	`products`.`category_layer3` AS `category_layer3`,
	`products`.`category_layer4` AS `category_layer4`,
	`products`.`title` AS `title`,
	`products`.`description` AS `description`,
	`products`.`style_title` AS `style_title`,
	`products`.`sub_style_title` AS `sub_style_title`,
	`products`.`style_enabled` AS `style_enabled`,
	`products`.`sub_style_enabled` AS `sub_style_enabled`,
	`products`.`min` AS `min`,
	`products`.`max` AS `max`,
	`products`.`photo` AS `photo`,
	`products`.`sold` AS `sold`,
	`products`.`is_public` AS `is_public`,
	`products`.`created_at` AS `created_at`,
	`products`.`updated_at` AS `updated_at`,
	`products`.`deleted_at` AS `deleted_at` 
FROM
	`products` 
WHERE
	`products`.`id` IN (
	SELECT
		`product_style_table`.`product_id` 
	FROM
		`product_style_table` 
	WHERE
		( `product_style_table`.`qty` <= `product_style_table`.`low_stock` ) 
GROUP BY
	`product_style_table`.`product_id`);

CREATE VIEW `report_over_sale` AS SELECT
	`products`.`id` AS `id`,
	`products`.`platform_id` AS `platform_id`,
	`products`.`category_layer1` AS `category_layer1`,
	`products`.`category_layer2` AS `category_layer2`,
	`products`.`category_layer3` AS `category_layer3`,
	`products`.`category_layer4` AS `category_layer4`,
	`products`.`title` AS `title`,
	`products`.`description` AS `description`,
	`products`.`style_title` AS `style_title`,
	`products`.`sub_style_title` AS `sub_style_title`,
	`products`.`style_enabled` AS `style_enabled`,
	`products`.`sub_style_enabled` AS `sub_style_enabled`,
	`products`.`min` AS `min`,
	`products`.`max` AS `max`,
	`products`.`photo` AS `photo`,
	`products`.`sold` AS `sold`,
	`products`.`is_public` AS `is_public`,
	`products`.`created_at` AS `created_at`,
	`products`.`updated_at` AS `updated_at`
FROM
	`products` 
WHERE
	`products`.`id` IN (
	SELECT
		`product_style_table`.`product_id` 
	FROM
		`product_style_table` 
	WHERE
		( `product_style_table`.`qty` < 0 ) 
GROUP BY
	`product_style_table`.`product_id`);

CREATE VIEW `report_wait_delivery_product` AS select products.* from orders
	inner join order_products as op on op.order_id = orders.id
	inner join products on products.id = op.product_id
	where orders.`status` = 21
	group by products.id;

CREATE VIEW `report_wait_delivery_style_table` AS select pst.*, count(pst.id) as wait_for_delivery from product_style_table as pst
inner join order_products as op on op.style_id = pst.id
inner join orders on orders.id = op.order_id
where orders.`status` = 21
group by pst.id;

ALTER TABLE `platform` 
ADD COLUMN `fb_pixel_token` varchar(255) NULL COMMENT 'conversion api token' AFTER `fb_pixel`;


CREATE TABLE `provider`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `platform_id` int NULL,
  `title` varchar(255) NULL,
  `job_title` varchar(255) NULL,
  `name` varchar(255) NULL,
  `mobile` varchar(255) NULL,
  `email` varchar(255) NULL,
  `line` varchar(255) NULL,
  `address` varchar(255) NULL,
  `bank` varchar(255) NULL,
  `account` varchar(255) NULL,
  `tax_id` varchar(8) NULL,
  `memo` tinytext NULL,
  `created_at` int(255) UNSIGNED NULL,
  `updated_at` int(255) UNSIGNED NULL,
  `deleted_at` int(255) UNSIGNED NULL,
  PRIMARY KEY (`id`)
);