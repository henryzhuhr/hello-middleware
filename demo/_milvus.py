"""
doc: https://milvus.io/docs/zh
"""

import json
import os
import random
import time
import argparse
from urllib import request
from pymilvus import (
    connections,
    Collection,
    MilvusClient,
    FieldSchema,
    CollectionSchema,
    DataType,
)
from pymilvus.client.types import (
    IndexType,
    MetricType,
)
from concurrent.futures import ThreadPoolExecutor, as_completed


def get_args():

    parser = argparse.ArgumentParser(description="Demo Milvus client")
    parser.add_argument(
        "--milvus-host",
        type=str,
        default=os.getenv("MILVUS_HOST", "milvus-standalone"),
        help="Milvus server host",
    )
    parser.add_argument(
        "--milvus-port",
        type=str,
        default=os.getenv("MILVUS_PORT", "19530"),
        help="Milvus server port",
    )
    parser.add_argument(
        "--milvus-user",
        type=str,
        default=os.getenv("MILVUS_USER", "root"),
        help="Milvus server user",
    )
    parser.add_argument(
        "--milvus-password",
        type=str,
        default=os.getenv("MILVUS_PASSWORD", "Milvus"),
        help="Milvus server password",
    )
    parser.add_argument("--collection-name", type=str, default="demo_collection")
    return parser.parse_args()


def main():
    args = get_args()
    print(args.__dict__)

    # 配置并连接到Milvus服务
    connections.connect(
        host=args.milvus_host,
        port=args.milvus_port,
        user=args.milvus_user,
        password=args.milvus_password,
    )

    # 创建Milvus客户端实例
    client = MilvusClient(
        uri=f"http://{args.milvus_host}:{args.milvus_port}",
        user=args.milvus_user,
        password=args.milvus_password,
        timeout=30,
    )
    print(
        "Successfully connected to Milvus:",
        {
            "server_version": client.get_server_version(),
        },
    )

    collections = client.list_collections()
    print("list_collections:", collections)

    collection_name: str = args.collection_name
    if collection_name in collections:
        # 删除集合
        client.drop_collection(collection_name)
        print(f"Collection '{collection_name}' deleted.")

    # 创建集合
    fields = [
        FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True),
        FieldSchema(name="text", dtype=DataType.VARCHAR, max_length=1024),
        FieldSchema(name="embeddings", dtype=DataType.FLOAT_VECTOR, dim=1024),
        FieldSchema(name="answer", dtype=DataType.VARCHAR, max_length=10000),
    ]
    schema = CollectionSchema(fields, description="no api collection")

    # 直接使用Collection类创建集合
    collection = Collection(name=collection_name, schema=schema)

    # 为embeddings字段创建索引
    collection.create_index(
        field_name="embeddings",
        index_params={
            # "index_type": IndexType.IVF_FLAT,
            "index_type": "IVF_FLAT",
            "params": {"nlist": 4096},
            "metric_type": MetricType.IP,
        },
    )
    collection.load()  # Load the data into memory.

    try:
        id = 0
        for batch in range(6):
            data = []
            for i in range(2):
                data.append(
                    {
                        "text": f"text_{id}",
                        "embeddings": [random.random() for i in range(1024)],
                        "answer": f"answer_{id}",
                    }
                )
                id += 1
            st = time.time()
            client.insert(collection_name, data)
            et = time.time()
            print(f"Batch {batch} inserted.", et - st)

        # 如果你知道集合中有多少条数据，可以设置limit参数来控制返回的数据量
        # 注意：如果集合非常大，这可能会导致内存问题
        results = client.query(
            collection_name,
            filter="id >= 0",
            output_fields=["id"],
        )
        print(f"Total records: {len(results)}")

    except Exception as e:
        print(e)
    finally:
        if collection:
            collection.release()  # Releases the collection data from memory.


if __name__ == "__main__":
    main()
