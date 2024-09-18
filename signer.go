package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

// сюда писать код

func ExecutePipeline(jobs ...job) {
	var ch chan interface{}

	for _, job := range jobs[:len(jobs)-1] {
		job := job
		in, out := ch, make(chan interface{})
		ch = out
		go func() {
			defer close(out)
			job(in, out)
		}()
	}

	jobs[len(jobs)-1](ch, nil)
}

func SingleHash(in, out chan interface{}) {
	type item struct {
		data string
		md5  string
	}

	var wg sync.WaitGroup
	ch := make(chan item)

	step2 := func() {
		defer wg.Done()
		v := <-ch
		a := make(chan string, 1)
		b := make(chan string, 1)
		go func() { a <- DataSignerCrc32(v.data) }()
		go func() { b <- DataSignerCrc32(v.md5) }()
		out <- <-a + "~" + <-b
	}

	for v := range in {
		var data string
		switch x := v.(type) {
		case int:
			data = strconv.Itoa(x)
		case string:
			data = x
		}

		wg.Add(1)
		go step2()
		ch <- item{data, DataSignerMd5(data)}
	}

	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	var wg sync.WaitGroup

	for v := range in {
		data := v.(string)

		wg.Add(1)
		go func() {
			defer wg.Done()

			ss := make([]string, 6) // 0..5
			var wg2 sync.WaitGroup
			wg2.Add(len(ss))

			for i := range ss {
				i := i
				go func() {
					defer wg2.Done()
					ss[i] = DataSignerCrc32(strconv.Itoa(i) + data)
				}()
			}

			wg2.Wait()
			out <- strings.Join(ss, "")
		}()
	}

	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var res []string

	for v := range in {
		data := v.(string)
		res = append(res, data)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})

	out <- strings.Join(res, "_")
}
