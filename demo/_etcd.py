import os
from typing import Tuple
import etcd3
from etcd3.client import KVMetadata


def get_args():
    import argparse

    parser = argparse.ArgumentParser(description="Demo etcd client")
    parser.add_argument(
        "--host",
        type=str,
        default=os.getenv("ETCD_HOST", "etcd-server"),
        help="Etcd server host",
    )
    parser.add_argument(
        "--port",
        type=int,
        default=os.getenv("ETCD_PORT", 2379),
        help="Etcd server port",
    )
    return parser.parse_args()


def main():
    args = get_args()
    print(args.__dict__)

    client = etcd3.client(host=args.host, port=args.port)

    # 设置一个键值对
    key = "demo_key"
    value = "demo_value"
    client.put(key, value)
    print(f"[Set] k-v {key}:{value}")

    # 从etcd获取值
    get_value: bytes
    get_value, _ = client.get(key)
    print(f"[Get] k-v {key}:{get_value.decode()}")

    # 列出某个目录下的所有键值对（如果适用）
    # 这里假设我们想要列出所有以'demo_'为前缀的键
    for event in client.get_prefix("demo_"):
        event: Tuple[bytes, KVMetadata]
        print(f"Found key: {event[1].key.decode()} with value: {event[0].decode()}")

    # 删除一个键
    is_del = client.delete(key)
    if is_del:
        print(f"Succeed to delete key {key}")
    else:
        print(f"Failed  to delete key {key}")

    # 尝试再次获取已删除的键的值
    get_value, _ = client.get(key)
    if get_value is None:
        print(f"Key {key} not found after deletion")
    else:
        print(f"Unexpectedly found value for key {key}: {get_value.decode()}")


if __name__ == "__main__":
    main()
