import json
import os
import traceback

from kafka import KafkaConsumer, KafkaProducer
from kafka.errors import KafkaError


class KafkaClient:
    def __init__(self, bootstrap_servers):
        self.bootstrap_servers = bootstrap_servers
        self.producer = KafkaProducer(
            bootstrap_servers=self.bootstrap_servers,
            allow_auto_create_topics=True,
            value_serializer=lambda v: json.dumps(v).encode("utf-8"),
            api_version=(4, 0, 0),
        )
        self.consumer = KafkaConsumer(
            bootstrap_servers=self.bootstrap_servers,
            value_deserializer=lambda x: json.loads(x.decode("utf-8")),
            api_version=(4, 0, 0),
        )

    def produce(self, topic, message):
        try:
            print(f"Producing message to topic {topic}: {message}")
            self.producer.send(topic, message)
            self.producer.flush(timeout=5)
            # self.producer.close()
        except KafkaError as e:
            print(f"Error producing message: {e}")
            traceback.print_exc()

    def consume(self, topic):
        try:
            self.consumer.subscribe([topic])
            for message in self.consumer:
                print(f"Received message: {message.value}")
        except KafkaError as e:
            print(f"Error consuming message: {e}")
            traceback.print_exc()


def main():
    bootstrap_servers = os.getenv("KAFKA_BROKER_HOST", "localhost")
    print(f"Connecting to Kafka broker at {bootstrap_servers}:9092")
    kafka_client = KafkaClient(bootstrap_servers=[f"{bootstrap_servers}:9092"])
    kafka_client.produce("test-topic", {"key": "value"})


if __name__ == "__main__":
    main()
