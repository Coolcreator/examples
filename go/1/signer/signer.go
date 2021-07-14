package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	out := make(chan interface{})
	var wgExecute sync.WaitGroup
	for _, j := range jobs {
		wgExecute.Add(1)
		go func(in, out chan interface{}, j job, wg *sync.WaitGroup) {
			j(in, out)
			close(out)
			wg.Done()
		}(in, out, j, &wgExecute)
		in = out
		out = make(chan interface{})
	}
	wgExecute.Wait()
}

func SingleHash(in, out chan interface{}) {
	var mutex sync.Mutex
	var wgSingle sync.WaitGroup
	for data := range in {
		wgSingle.Add(1)
		go func(input interface{}, wgSingle *sync.WaitGroup) {
			data := strconv.Itoa(input.(int))
			hashes := []string{data}
			mutex.Lock()
			hashes = append(hashes, DataSignerMd5(data))
			mutex.Unlock()
			var wgSingleCrc32 sync.WaitGroup
			for index, _ := range hashes {
				wgSingleCrc32.Add(1)
				go func(i int, wgSingleCrc32 *sync.WaitGroup) {
					hashes[i] = DataSignerCrc32(hashes[i])
					wgSingleCrc32.Done()
				}(index, &wgSingleCrc32)
			}
			wgSingleCrc32.Wait()
			out <- strings.Join(hashes, "~")
			wgSingle.Done()
		}(data, &wgSingle)
	}
	wgSingle.Wait()
}

func MultiHash(in, out chan interface{}) {
	var wgMulti sync.WaitGroup
	for data := range in {
		wgMulti.Add(1)
		go func(input interface{}, wgMulti *sync.WaitGroup) {
			data := input.(string)
			hashes := []string{"0", "1", "2", "3", "4", "5"}
			var wgMultiCrc32 sync.WaitGroup
			for index, _ := range hashes {
				wgMultiCrc32.Add(1)
				go func(i int, wgMultiCrc32 *sync.WaitGroup) {
					hashes[i] = DataSignerCrc32(hashes[i] + data)
					wgMultiCrc32.Done()
				}(index, &wgMultiCrc32)
			}
			wgMultiCrc32.Wait()
			out <- strings.Join(hashes, "")
			wgMulti.Done()
		}(data, &wgMulti)
	}
	wgMulti.Wait()
}

func CombineResults(in, out chan interface{}) {
	hashes := []string{}
	for data := range in {
		hashes = append(hashes, data.(string))
	}
	sort.Strings(hashes)
	out <- strings.Join(hashes, "_")
}
