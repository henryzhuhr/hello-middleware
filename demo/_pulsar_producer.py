from datetime import datetime
import os

import pulsar
from loguru import logger


def main():
    """
    生产者程序，用于向指定的主题发送消息。
    """
    # 配置 Pulsar Broker 的服务 URL 和测试主题
    broker_url = os.getenv("PULSAR_BROKER_URL", "pulsar://localhost:6650")
    topic = os.getenv("PULSAR_TOPIC", "persistent://public/default/test-topic")

    # 执行生产者逻辑
    run_producer(broker_url, topic)


def run_producer(broker_url: str, topic: str):
    """
    运行生产者逻辑。
    :param broker_url: Pulsar Broker 的服务 URL (例如: pulsar://localhost:6650)
    :param topic: 测试的主题名称 (例如: persistent://public/default/test-topic)
    """
    try:
        client = pulsar.Client(broker_url)

        # 创建生产者
        producer = client.create_producer(topic)
        logger.info(f"Successfully created producer, topic: {topic}")

        # 发送测试消息
        test_message = f"Hello, Pulsar at {datetime.now()}"
        producer.send(test_message.encode("utf-8"))
        logger.info(f"Message sent: {test_message}")

        # 关闭客户端
        client.close()

    except Exception as e:
        logger.info(f"Connection or operation failed: {e}")


if __name__ == "__main__":
    main()


# curl http://pulsar-cluster-broker:8080/admin/v2/brokers/health