import argparse
import socket
import time
import threading


def client_func(client):
    url = "https://www.lipsum.com/feed/html"
    print('URL: ', url)
    client.sendall(bytes(url, 'UTF-8'))
    in_data = client.recv(1024)
    print(f"From Server with {url}:" , in_data.decode(), type(in_data))
    # out_data = 'bye'
    # client.sendall(bytes(out_data,'UTF-8'))
    client.close()


def run():
    parser = argparse.ArgumentParser()
    parser.add_argument('amount_of_threads', type=int, help='Amount of threads')
    parser.add_argument('file_path', type=str, help='urls file path')
    args = parser.parse_args()
    print(args.amount_of_threads, args.file_path)

    SERVER = "127.0.0.1"
    PORT = 8080

    t1 = time.time()

    client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client.connect((SERVER, PORT))
    t = threading.Thread(target = client_func, args = (client,))
    t.start()
    t.join()

    t2 = time.time()
    print('single time: ', t2 - t1)

    t1 = time.time()
    threads = []

    for i in range(args.amount_of_threads):
        client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        client.connect((SERVER, PORT))
        t = threading.Thread(target = client_func, args = (client,))
        t.start()
        threads.append(t)

    for x in threads:
        x.join()

    t2 = time.time()
    print('thread time: ', t2 - t1)


if __name__ == "__main__":
    run()
