import os

from demo._kafka_producer import KafkaClient


def main():
    bootstrap_servers = os.getenv("KAFKA_BROKER_HOST", "localhost")
    print(f"Connecting to Kafka broker at {bootstrap_servers}:9092")
    kafka_client = KafkaClient(bootstrap_servers=[f"{bootstrap_servers}:9092"])
    kafka_client.consume("test-topic")


if __name__ == "__main__":
    main()
