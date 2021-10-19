ALTER TABLE `ec`.`page_component_data_draft` 
ADD COLUMN `link_type` tinyint(4) NULL COMMENT '0: 無, 1: 頁面, 2: 分類, 3: 產品, 4: 外部連結' AFTER `img`;

ALTER TABLE `ec`.`page_component_data` 
ADD COLUMN `link_type` tinyint(4) NULL COMMENT '0: 無, 1: 頁面, 2: 分類, 3: 產品, 4: 外部連結' AFTER `img`;