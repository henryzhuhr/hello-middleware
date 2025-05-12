import os

from pymongo import MongoClient


def main():
    # 替换为你的 MongoDB 认证信息
    username = os.getenv("MONGO_INITDB_ROOT_USERNAME", "root")
    password = os.getenv("MONGO_INITDB_ROOT_PASSWORD", "root")
    host = os.getenv("MONGO_HOST", "localhost")
    port = 27017
    auth_db = "admin"  # 通常是 admin 数据库，取决于你的 MongoDB 配置

    # 构造带认证的连接字符串
    connection_string = f"mongodb://{username}:{password}@{host}:{port}/{auth_db}"

    # 连接 MongoDB
    client = MongoClient(connection_string)

    # Create or switch to a database
    db = client["mydatabase"]

    # Create or switch to a collection
    collection = db["mycollection"]

    # Insert a document
    document = {"name": "Tim", "age": 30}
    collection.insert_one(document)

    # Query the collection
    result = collection.find_one({"name": "Tim"})
    print(result)


if __name__ == "__main__":
    main()
