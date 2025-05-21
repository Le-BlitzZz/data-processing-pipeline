BROKER_SERVER = "rabbitmq:5672"
BROKER_USER = "etlstream"
BROKER_PASSWORD = "etlstream"

BROKER_RAW_PROCESSING_EXCHANGE = "raw_processing_exchange"
BROKER_PROCESSED_EXCHANGE = "processed_exchange"

BROKER_RAW_PROCESSING_QUEUE = "raw_processing_queue"

BATCH_SIZE = 11_986


def broker_dsn():
    return f"amqp://{BROKER_USER}:{BROKER_PASSWORD}@{BROKER_SERVER}/"
