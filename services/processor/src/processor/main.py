import processor.processor as processor
import processor.config as config
import asyncio


async def _app():
    broker = await config.new_broker()

    try:
        await processor.run(broker)
    finally:
        await config.shutdown(broker)


def main():
    config.setup_logging()

    asyncio.run(_app())


if __name__ == "__main__":
    main()
