ALTER TABLE `ec`.`page_component_data_draft` 
ADD COLUMN `link_type` tinyint(4) NULL COMMENT '0: 無, 1: 頁面, 2: 分類, 3: 產品, 4: 外部連結' AFTER `img`,
ADD COLUMN `link_id` int NULL COMMENT '連結的頁面/分類/產品 id' AFTER `link_type`;

ALTER TABLE `ec`.`page_component_data` 
ADD COLUMN `link_type` tinyint(4) NULL COMMENT '0: 無, 1: 頁面, 2: 分類, 3: 產品, 4: 外部連結' AFTER `img`,
ADD COLUMN `link_id` int NULL COMMENT '連結的頁面/分類/產品 id' AFTER `link_type`;