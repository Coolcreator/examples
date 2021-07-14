import argparse
import requests
import socket
import threading
import json
from collections import Counter

requests_counter = 0

def download(url, lock, k_freq_words):
    global requests_counter  #, lock

    resp = requests.get(url)
    data = resp.content.decode("utf-8")
    split_it = data.split()
    Counters_found = Counter(split_it)
    most_occur = Counters_found.most_common(k_freq_words)

    lock.acquire() 
    requests_counter += 1
    lock.release()
    print('urls processed: ', requests_counter)
    print(most_occur, type(most_occur), json.dumps(most_occur))

    return most_occur


def worker(server, lock, k_freq_words, i):
    while True:
        print('before accept connection', server, i)
        clientsock, clientAddress = server.accept()
        url = clientsock.recv(2048)
        print(url, clientAddress)
        n_most_freq_words = download(url, lock, k_freq_words)
        clientsock.send(bytes(json.dumps({k:v for k, v in n_most_freq_words}), 'UTF-8'))
        print ("Client disconnected...")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('-w', type=int, help='Number of workers')
    parser.add_argument('-k', type=int, help='Top frequent words')
    args = parser.parse_args()
    k_freq_words = args.k

    lock = threading.Lock()

    LOCALHOST = "127.0.0.1"
    PORT = 8080

    server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)  # lock передать в worker
    server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    server.bind((LOCALHOST, PORT))
    print("Server started")
    print("Waiting for client request..")
    server.listen()
    for i in range(args.w):
        t = threading.Thread(target = worker, args = (server, lock, k_freq_words, i))
        t.start()


if __name__ == "__main__":
    main()
