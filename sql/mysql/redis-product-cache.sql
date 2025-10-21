-- 创建数据库
CREATE DATABASE IF NOT EXISTS redis_product_cache CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

-- 使用数据库
USE redis_product_cache;


CREATE TABLE IF NOT EXISTS `products` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '商品ID',
  `name` VARCHAR(255) NOT NULL COMMENT '商品名称',
  `description` TEXT COMMENT '商品详情（富文本）',
  `category_id` BIGINT NOT NULL COMMENT '类目ID',
  `brand` VARCHAR(100) DEFAULT NULL COMMENT '品牌',
  `price` DECIMAL(10,2) NOT NULL COMMENT '原价（单位：元）',
  `sale_price` DECIMAL(10,2) DEFAULT NULL COMMENT '促销价（单位：元）',
  `main_image` VARCHAR(500) DEFAULT NULL COMMENT '主图URL',
  `images` JSON DEFAULT NULL COMMENT '商品轮播图URL列表',
  `attributes` JSON DEFAULT NULL COMMENT '自定义属性，如 {"材质": "纯棉", "适用季节": "夏季"}',
  `weight` DECIMAL(8,3) DEFAULT 0.000 COMMENT '重量（kg）',
  `free_shipping` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否包邮：0-否，1-是',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态：active, inactive, deleted',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_brand` (`brand`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='商品主表（SPU）';

CREATE TABLE IF NOT EXISTS `product_skus`(
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'SKU ID',
  `product_id` BIGINT NOT NULL COMMENT '关联商品ID',
  `specs` JSON NOT NULL COMMENT '规格组合，如 {"颜色": "红色", "尺寸": "L"}',
  `price` DECIMAL(10,2) NOT NULL COMMENT '该SKU售价',
  `stock` BIGINT NOT NULL DEFAULT 0 COMMENT '库存',
  `sku_code` VARCHAR(100) DEFAULT NULL COMMENT '商家自定义SKU编码',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sku_code` (`sku_code`),  -- 若启用SKU编码，需唯一
  KEY `idx_product_id` (`product_id`),
  KEY `idx_stock` (`stock`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='商品SKU表';


-- CREATE TABLE `categories` (
--   `id` BIGINT NOT NULL AUTO_INCREMENT,
--   `name` VARCHAR(100) NOT NULL,
--   `parent_id` BIGINT DEFAULT 0 COMMENT '父类目ID，0表示根类目',
--   `level` TINYINT NOT NULL DEFAULT 1 COMMENT '层级：1-一级，2-二级...',
--   `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  
--   PRIMARY KEY (`id`),
--   KEY `idx_parent_id` (`parent_id`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;