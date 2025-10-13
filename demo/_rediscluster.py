import argparse
import os

import redis


def get_args():
    parser = argparse.ArgumentParser(description="Demo Redis client")
    parser.add_argument(
        "--host",
        type=str,
        default=os.getenv("REDIS_HOST", "redis-server"),
        help="Redis server host",
    )
    parser.add_argument(
        "--port",
        type=int,
        default=os.getenv("REDIS_PORT", 6379),
        help="Redis server port",
    )
    return parser.parse_args()


def main():
    args = get_args()
    client = redis.StrictRedis(
        host=args.host,
        port=args.port,
        db=0,
        decode_responses=True,
    )

    # 设置一个键值对
    key = "demo_key"
    value = "demo_value"
    client.set(key, value)
    print(f"Set key: {key} with value: {value}")

    # 从Redis获取值
    get_value = client.get(key)
    print(f"Get value for key {key}: {get_value}")

    # 检查键是否存在
    exists = client.exists(key)
    print(f"Does key {key} exist? {'Yes' if exists else 'No'}")

    # 使用列表 - 向列表尾部添加元素
    list_key = "demo_list"
    elements = ["element1", "element2", "element3"]
    for element in elements:
        client.rpush(list_key, element)
    print(f"Added elements to list {list_key}")

    # 获取列表中的所有元素
    list_elements = client.lrange(list_key, 0, -1)
    print(f"List {list_key} contents: {list_elements}")

    # 删除一个键
    client.delete(key)
    print(f"Deleted key: {key}")

    # 尝试再次获取已删除的键的值
    deleted_value = client.get(key)
    if deleted_value is None:
        print(f"Key {key} not found after deletion")
    else:
        print(f"Unexpectedly found value for key {key}: {deleted_value}")


if __name__ == "__main__":
    main()
