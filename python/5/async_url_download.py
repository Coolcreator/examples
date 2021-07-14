import aiohttp
import argparse
import asyncio
import time
import sys


parser = argparse.ArgumentParser()
parser.add_argument('-c', type=int,
                    help='number of threads')
parser.add_argument('file_path', type=str,
                    help='file with urls')
args = parser.parse_args()

sem = asyncio.Semaphore(args.c)

async def fetch_url(session, url):
    async with session.get(url) as resp:
        pokemon = await resp.json()
        return pokemon


async def sem_fetcher(session, url):
    async with sem:  # semaphore limits num of simultaneous downloads
        return await fetch_url(session, url)


async def main():
    async with aiohttp.ClientSession() as session:
        tasks = []
        with open(args.file_path) as file:
            for url in file:
                tasks.append(asyncio.create_task(sem_fetcher(session, url)))

        results = await asyncio.gather(*tasks)
        for result in results:
            print(result)

start_time = time.time()

if __name__ ==  '__main__':
    loop = asyncio.get_event_loop()
    try:
        loop.run_until_complete(main())
    finally:
        loop.run_until_complete(loop.shutdown_asyncgens())
        loop.close()
        
print(f"Executed in {time.time() - start_time} seconds.")
