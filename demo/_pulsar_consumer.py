import os

import pulsar
from loguru import logger


def main():
    """
    消费者程序，用于从指定的主题接收消息。
    """
    # 配置 Pulsar Broker 的服务 URL 和测试主题
    broker_url = os.getenv("PULSAR_BROKER_URL", "pulsar://localhost:6650")
    topic = os.getenv("PULSAR_DEFAULT_TOPIC", "persistent://public/default/test-topic")

    # 执行消费者逻辑
    run_consumer(broker_url, topic)


def run_consumer(broker_url: str, topic: str):
    """
    运行消费者逻辑。
    :param broker_url: Pulsar Broker 的服务 URL (例如: pulsar://localhost:6650)
    :param topic: 测试的主题名称 (例如: persistent://public/default/test-topic)
    """
    try:
        client = pulsar.Client(broker_url)

        # 创建消费者
        consumer = client.subscribe(topic, subscription_name="test-subscription")
        logger.info(f"Successfully created consumer, subscribed to theme: {topic}")

        while True:
            # 接收消息
            msg = consumer.receive(timeout_millis=10000)  # 设置超时时间为 10 秒
            received_message = msg.data().decode("utf-8")
            logger.info(f"receive message: {received_message}")

            # 确认消息
            consumer.acknowledge(msg)
            logger.info("Message confirmed")

    except Exception as e:
        logger.info(f"Connection or operation failed: {e}")


if __name__ == "__main__":
    main()
