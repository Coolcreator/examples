#include <aio.h>
#include <errno.h>
#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include  "main.h"

#define nbytes 1024

char* async_read(char* filename) {
    int fd = open(filename, O_RDONLY);
    if (fd == -1) {
        printf("Can't open file %s!\n", filename);
        exit(EXIT_FAILURE);
    }

    struct aiocb cb;
    int bs = nbytes;
    char* buffer = (char*)malloc(bs * sizeof(char));
    if (buffer == NULL) {
        printf("Malloc error!\n");
        exit(EXIT_FAILURE);
    }

    memset(&cb, 0, sizeof(struct aiocb));
    cb.aio_fildes = fd;
    cb.aio_offset = 0;
    cb.aio_nbytes = nbytes;
    cb.aio_buf = buffer;

    while (1) {
        yield();
        if (aio_read(&cb) == -1) {
            printf("Unable to create request!\n");
            free(buffer);
            close(fd);
            exit(EXIT_FAILURE);
        }
        yield();
        while (aio_error(&cb) == EINPROGRESS);
        yield();
        int nb = aio_return(&cb);
        if (nb == -1) {
            printf("Error happened while aio_return!\n");
            free(buffer);
            close(fd);
            exit(EXIT_FAILURE);
        }
        yield();
        cb.aio_offset += nb;
        yield();
        if (bs % cb.aio_offset == 0) {
            yield();
            bs = cb.aio_offset + nbytes;
            yield();
            buffer = (char*)realloc(buffer, bs);
            yield();
            if (buffer == NULL) {
                printf("Realloc error!\n");
                exit(EXIT_FAILURE);
            }
        } else {
            yield();
            bs = cb.aio_offset + 1;
            yield();
            buffer = (char*)realloc(buffer, bs);
            if (buffer == NULL) {
                printf("Realloc error!\n");
                exit(EXIT_FAILURE);
            }
            yield();
            buffer[bs-1] = '\0';
            yield();
            close(fd);
            return buffer;
        }
        cb.aio_buf = buffer + cb.aio_offset;
    }
}