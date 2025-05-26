import logging
from aio_pika import connect_robust

BROKER_SERVER = "rabbitmq:5672"
BROKER_USER = "etlstream"
BROKER_PASSWORD = "etlstream"

BROKER_RAW_PROCESSING_EXCHANGE = "raw_processing_exchange"
BROKER_PROCESSED_EXCHANGE = "processed_exchange"

BROKER_RAW_PROCESSING_QUEUE = "raw_processing_queue"

BATCH_SIZE = 11_986


async def new_broker():
    broker = await connect_robust(broker_dsn())
    logging.info("RabbitMQ: connection established")
    return broker


async def shutdown(broker):
    await broker.close()
    logging.info("closed message broker connection")


def setup_logging(level=logging.INFO):
    logging.basicConfig(
        level=level,
        format="%(asctime)s %(levelname)s %(message)s",
        datefmt="%H:%M:%S",
    )


def broker_dsn():
    return f"amqp://{BROKER_USER}:{BROKER_PASSWORD}@{BROKER_SERVER}/"
