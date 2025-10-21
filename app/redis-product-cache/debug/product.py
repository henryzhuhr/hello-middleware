# client.py
import json
import time
from typing import Any, Dict, List, Optional

import requests

BASE_URL = "http://localhost:8888/v1"


class ProductClient:
    def __init__(self, base_url: str = BASE_URL):
        self.base_url = base_url.rstrip("/")
        self.session = requests.Session()
        self.session.headers.update({"Content-Type": "application/json"})

    def add_product(
        self,
        name: str,
        description: str,
        category_id: int,
        price: float,
        brand: Optional[str] = None,
        sale_price: Optional[float] = None,  # 可选促销价
        attributes: Optional[Dict[str, str]] = None,
        skus: Optional[List[Dict[str, Any]]] = None,
    ) -> int:
        """
        添加商品
        :return: 商品 ID
        """
        payload = {
            "name": name,
            "description": description,
            "categoryId": category_id,
            "price": price,
        }

        # 可选字段：只有非 None 时才加入
        if brand is not None:
            payload["brand"] = brand
        if sale_price is not None:
            payload["salePrice"] = sale_price
        if attributes is not None:
            payload["attributes"] = attributes
        if skus is not None:
            payload["skus"] = skus

        url = f"{self.base_url}/product"
        response = self.session.post(url, data=json.dumps(payload))

        try:
            response.raise_for_status()
            resp_json = response.json()
            return resp_json["id"]
        except requests.exceptions.HTTPError as e:
            print(f"HTTP Error: {e}")
            print("Response:", response.text)
            raise
        except (KeyError, json.JSONDecodeError) as e:
            print(f"Invalid response format: {e}")
            print("Raw response:", response.text)
            raise

    def get_product(self, product_id: int) -> dict:
        """
        查询商品详情
        :return: 商品详情字典
        """
        url = f"{self.base_url}/product/{product_id}"
        response = self.session.get(url)

        try:
            response.raise_for_status()
            return response.json()
        except requests.exceptions.HTTPError as e:
            print(f"HTTP Error: {e}")
            print("Response:", response.text)
            raise
        except json.JSONDecodeError as e:
            print(f"JSON Decode Error: {e}")
            print("Raw response:", response.text)
            raise


def main():
    client = ProductClient()

    # --- 示例 1：添加一个无促销价的商品 ---
    print("📦 添加普通商品...")
    product_id = client.add_product(
        name="纯棉T恤",
        description="100%纯棉，舒适透气",
        category_id=101,
        price=99.9,
        brand="MyBrand",
        attributes={"材质": "纯棉", "季节": "夏季"},
        skus=[
            {
                "specs": {"颜色": "白色", "尺码": "M"},
                "price": 99.9,
                "stock": 100,
                "skuCode": "TSHIRT-WHT-M",
            },
            {
                "specs": {"颜色": "黑色", "尺码": "L"},
                "price": 99.9,
                "stock": 50,
                "skuCode": "TSHIRT-BLK-L",
            },
        ],
    )
    print(f"✅ 商品创建成功，ID: {product_id}")

    # --- 示例 2：添加一个促销价为 0 的免费商品 ---
    print("\n🎁 添加免费赠品...")
    free_product_id = client.add_product(
        name="会员专属赠品",
        description="仅限VIP用户领取",
        category_id=200,
        price=10.0,  # 原价
        sale_price=0.0,  # 促销价为 0 元（免费）
        skus=[
            {
                "specs": {"类型": "电子券"},
                "price": 0.0,
                "stock": 1000,
                "skuCode": "GIFT-VIP-001",
            }
        ],
    )
    print(f"✅ 免费商品创建成功，ID: {free_product_id}")

    # --- 示例 3：查询商品 ---
    print(f"\n🔍 查询商品 ID={product_id}...")
    product_detail = client.get_product(product_id)
    print("商品名称:", product_detail["name"])
    print("原价:", product_detail["price"])
    print("促销价:", product_detail.get("salePrice", "未设置"))  # 可能不存在
    print("总库存:", product_detail["totalStock"])
    print("SKU 数量:", len(product_detail["skus"]))

    # --- 示例 4：查询免费商品 ---
    print(f"\n🔍 查询免费商品 ID={free_product_id}...")
    free_detail = client.get_product(free_product_id)
    print("促销价:", free_detail.get("salePrice"))  # 应输出 0.0


if __name__ == "__main__":
    main()
