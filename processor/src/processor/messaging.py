import asyncio, logging, os
from aio_pika import (
    connect_robust,
    Message,
    DeliveryMode,
    ExchangeType,
)

import processor.config as config
from processor.transformation import transform

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(message)s",
    datefmt="%H:%M:%S",
)


async def main():
    logging.info("Processor: starting")

    conn = await connect_robust(config.broker_dsn())
    logging.info("RabbitMQ: connection established")

    channel = await conn.channel()
    await channel.set_qos(prefetch_count=os.cpu_count())

    raw_ex = await channel.declare_exchange(
        config.BROKER_RAW_PROCESSING_EXCHANGE, ExchangeType.FANOUT, durable=True
    )
    proc_ex = await channel.declare_exchange(
        config.BROKER_PROCESSED_EXCHANGE, ExchangeType.FANOUT, durable=True
    )

    queue = await channel.declare_queue(
        config.BROKER_RAW_PROCESSING_QUEUE, durable=True
    )
    await queue.bind(raw_ex)

    rows = []
    batch_ready = asyncio.Event()

    async def on_message(msg):
        async with msg.process(requeue=False):
            rows.append(msg.body)
            if len(rows) >= config.BATCH_SIZE:
                batch_ready.set()

    consumer_tag = await queue.consume(on_message)
    logging.info("Waiting for %d messages …", config.BATCH_SIZE)

    await batch_ready.wait()

    await _process_and_publish(proc_ex, rows)

    await queue.cancel(consumer_tag)
    await channel.close()
    await conn.close()

    logging.info("closed message broker connection")


async def _process_and_publish(proc_ex, raw_rows):
    logging.info("received batch - transforming with pandas")

    processed = await asyncio.to_thread(transform, raw_rows)

    for row in processed:
        await proc_ex.publish(
            Message(
                body=row.encode(),
                content_type="application/json",
                delivery_mode=DeliveryMode.PERSISTENT,
            ),
            routing_key="",
        )

    logging.info("Processor: batch processed and rows published")
