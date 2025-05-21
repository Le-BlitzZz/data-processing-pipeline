from processor.messaging import main as _run_async
import asyncio


def main():
    asyncio.run(_run_async())


if __name__ == "__main__":
    main()
