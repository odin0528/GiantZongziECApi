ALTER TABLE `ec`.`page_component_data_draft` 
ADD COLUMN `link_type` tinyint(4) NULL COMMENT '0: 無, 1: 頁面, 2: 分類, 3: 產品, 4: 外部連結' AFTER `img`;

ALTER TABLE `ec`.`page_component_data` 
ADD COLUMN `link_type` tinyint(4) NULL COMMENT '0: 無, 1: 頁面, 2: 分類, 3: 產品, 4: 外部連結' AFTER `img`;

CREATE TABLE `ec`.`platform_menu`  (
  `id` int NOT NULL,
  `platform_id` int NOT NULL,
  `title` varchar(255) NOT NULL,
  `link_type` tinyint(4) NULL,
  `link` varchar(255) NULL,
  `sort` tinyint(4) NULL,
  `is_enabled` tinyint(4) NULL,
  `created_at` int NULL,
  `updated_at` int NULL,
  `deleted_at` int NULL,
  PRIMARY KEY (`id`)
);

ALTER TABLE `ec`.`pages` 
DROP COLUMN `is_menu`,
DROP COLUMN `sort`;

ALTER TABLE `ec`.`platform_menu` 
MODIFY COLUMN `id` int(11) NOT NULL AUTO_INCREMENT FIRST;