import asyncio, logging, os, json
from aio_pika import (
    Message,
    DeliveryMode,
    ExchangeType,
)

import processor.config as config
from processor.transform import transform


async def run(conn):
    channel = await conn.channel()
    await channel.set_qos(prefetch_count=os.cpu_count())

    raw_processing_exchange = await channel.declare_exchange(
        config.BROKER_RAW_PROCESSING_EXCHANGE, ExchangeType.FANOUT, durable=True
    )
    processed_exchange = await channel.declare_exchange(
        config.BROKER_PROCESSED_EXCHANGE, ExchangeType.FANOUT, durable=True
    )

    queue = await channel.declare_queue(
        config.BROKER_RAW_PROCESSING_QUEUE, durable=True
    )
    await queue.bind(raw_processing_exchange)

    rows = []
    batch_ready = asyncio.Event()

    async def on_message(msg):
        async with msg.process(requeue=False):
            rows.append(msg.body)
            if len(rows) >= config.BATCH_SIZE:
                batch_ready.set()

    consumer_tag = await queue.consume(on_message)
    logging.info("Waiting for %d messages", config.BATCH_SIZE)

    await batch_ready.wait()

    await _process_and_publish(processed_exchange, rows)

    await queue.cancel(consumer_tag)
    await channel.close()


async def _process_and_publish(proc_ex, raw_rows):
    logging.info("received batch - transforming with pandas")

    json_rows = [json.loads(raw_row) for raw_row in raw_rows]
    processed = await asyncio.to_thread(transform, json_rows)

    for row in processed:
        await proc_ex.publish(
            Message(
                body=json.dumps(row).encode(),
                content_type="application/json",
                delivery_mode=DeliveryMode.PERSISTENT,
            ),
            routing_key="",
        )

    logging.info("Processor: batch processed and rows published")
