from elasticsearch import Elasticsearch, helpers
import os
from typing import Dict, Any


def get_args():
    import argparse

    parser = argparse.ArgumentParser(description="Demo Elasticsearch client")
    parser.add_argument(
        "--host",
        type=str,
        default=os.getenv("ELASTICSEARCH_HOST", "elasticsearch-server"),
        help="Elasticsearch server host",
    )
    parser.add_argument(
        "--port",
        type=int,
        default=os.getenv("ELASTICSEARCH_PORT", 9200),
        help="Elasticsearch server port",
    )
    return parser.parse_args()


def main():
    args = get_args()
    print(args.__dict__)

    # 创建Elasticsearch客户端实例，并指定scheme
    es = Elasticsearch([{"host": args.host, "port": args.port, "scheme": "http"}])

    # 检查Elasticsearch服务是否可用
    if not es.ping():
        print("Cannot connect to Elasticsearch!")
        return
    else:
        print("Connected to Elasticsearch successfully!")

    index_name = "demo_index"

    # 创建索引
    if not es.indices.exists(index=index_name):
        es.indices.create(index=index_name)
        print(f"Index '{index_name}' created.")
    else:
        print(f"Index '{index_name}' already exists.")

    # 索引文档
    doc: Dict[str, Any] = {
        "title": "Example document",
        "text": "This is a demo text.",
        "timestamp": "2025-02-13",
    }
    res = es.index(index=index_name, id=1, body=doc)
    print(f"Indexed with result: {res['result']}")

    # 刷新索引以确保文档可以被搜索到
    es.indices.refresh(index=index_name)

    # 搜索文档
    query_body = {"query": {"match": {"text": "demo"}}}
    res = es.search(index=index_name, body=query_body)
    print("Search results:")
    for hit in res["hits"]["hits"]:
        print(hit["_source"])


if __name__ == "__main__":
    main()
