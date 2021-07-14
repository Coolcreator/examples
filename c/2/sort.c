#include <stdio.h>

void merge(int arr[], int lb, int md, int rb) {
    int s1 = md - lb + 1;
    int s2 = rb - md;

    int larr[s1], rarr[s2];

    for (int i = 0; i < s1; i++) {
        larr[i] = arr[lb+i];
    }

    for (int j = 0; j < s2; j++) {
        rarr[j] = arr[md+j+1];
    }

    int i = 0, j = 0, k = lb;
    while (i < s1 && j < s2) {
        if (larr[i] <= rarr[j]) {
            arr[k++] = larr[i++];
        } else {
            arr[k++] = rarr[j++];
        }
    }

    while (i < s1) {
        arr[k++] = larr[i++];
    }

    while (j < s2) {
        arr[k++] = rarr[j++];        
    }
    
    return;
}

void merge_sort(int arr[], int lb, int rb) {
    if (lb < rb) {
        int md = lb + (rb-lb) / 2;

        merge_sort(arr, lb, md);
        merge_sort(arr, md+1, rb);
        merge(arr, lb, md, rb);
    }

    return;
}

void print_array(int arr[], int len) {
    for (int i = 0; i < len; i++) {
        printf("%d ", arr[i]);
    }
    printf("\n");

    return;
}
