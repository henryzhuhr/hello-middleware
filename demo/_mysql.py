import os
import argparse
import mysql.connector
from mysql.connector import Error


def get_args():
    parser = argparse.ArgumentParser(description="MySQL connection arguments")
    parser.add_argument(
        "--host", default=os.getenv("MYSQL_HOST", "mysql-server"), help="MySQL host"
    )
    parser.add_argument(
        "--database",
        default=os.getenv("MYSQL_DATABASE", "mysql"),
        help="MySQL database",
    )
    parser.add_argument(
        "--user", default=os.getenv("MYSQL_USER", "root"), help="MySQL user"
    )
    parser.add_argument(
        "--password",
        default=os.getenv("MYSQL_PASSWORD", "password"),
        help="MySQL password",
    )
    return parser.parse_args()


def main():
    args = get_args()
    print(args.__dict__)
    connection=None
    try:
        # 创建连接
        connection = mysql.connector.connect(
            host=args.host,
            database=args.database,
            user=args.user,
            password=args.password,
        )

        if connection.is_connected():
            db_Info = connection.get_server_info()
            print(f"Connected to MySQL Server version {db_Info}")

            cursor = connection.cursor()

            # 创建表（如果需要）
            cursor.execute(
                "CREATE TABLE IF NOT EXISTS demo_table (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255), age INT)"
            )
            print("Table created successfully")

            # 插入数据
            insert_query = """INSERT INTO demo_table (name, age) VALUES (%s, %s)"""
            values = ("Alice", 30)
            cursor.execute(insert_query, values)
            connection.commit()
            print(
                f"Record inserted successfully into demo_table with id {cursor.lastrowid}"
            )

            # 查询数据
            cursor.execute("SELECT * FROM demo_table")
            records = cursor.fetchall()
            for row in records:
                print(f"id={row[0]}, name={row[1]}, age={row[2]}")

            # 删除数据
            delete_query = "DELETE FROM demo_table WHERE name = 'Alice'"
            cursor.execute(delete_query)
            connection.commit()
            print(f"{cursor.rowcount} record(s) deleted")
    except Error as e:
        print(f"Error while connecting to MySQL: {e}")
    except Exception as e:
        print(f"Error while connecting to MySQL: {e}")
    finally:
        if connection and connection.is_connected():
            cursor.close()
            connection.close()
            print("MySQL connection is closed")


if __name__ == "__main__":
    main()
