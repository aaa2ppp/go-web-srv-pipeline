package main

import (
	"reflect"
	"testing"
)

func TestSingleHash(t *testing.T) {
	tests := []struct {
		name string
		data string
		want string
	}{
		{
			"0",
			"0",
			"4108050209~502633748",
		},
		{
			"1",
			"1",
			"2212294583~709660146",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan interface{})
			out := make(chan interface{})
			go SingleHash(in, out)
			in <- tt.data
			close(in)
			got := (<-out).(string)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMultiHash(t *testing.T) {
	tests := []struct {
		name string
		data string
		want string
	}{
		{
			"0",
			"4108050209~502633748",
			"29568666068035183841425683795340791879727309630931025356555",
		},
		{
			"1",
			"2212294583~709660146",
			"4958044192186797981418233587017209679042592862002427381542",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan interface{})
			out := make(chan interface{})
			go MultiHash(in, out)
			in <- tt.data
			close(in)
			got := (<-out).(string)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCombineResults(t *testing.T) {
	tests := []struct {
		name string
		data []string
		want string
	}{
		{
			"0",
			[]string{
				"29568666068035183841425683795340791879727309630931025356555",
				"4958044192186797981418233587017209679042592862002427381542",
			},
			"29568666068035183841425683795340791879727309630931025356555_4958044192186797981418233587017209679042592862002427381542",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan interface{})
			out := make(chan interface{})
			go CombineResults(in, out)
			for _, v := range tt.data {
				in <- v
			}
			close(in)
			got := (<-out).(string)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExecutePipeline(t *testing.T) {
	tests := []struct {
		name string
		jobs []job
		data []string
		want []string
	}{
		{
			"1",
			[]job{SingleHash},
			[]string{"0", "1"},
			[]string{"4108050209~502633748", "2212294583~709660146"},
		},
		{
			"2",
			[]job{SingleHash, MultiHash},
			[]string{"0", "1"},
			[]string{
				"29568666068035183841425683795340791879727309630931025356555",
				"4958044192186797981418233587017209679042592862002427381542",
			},
		},
		{
			"3",
			[]job{SingleHash, MultiHash, CombineResults},
			[]string{"0", "1"},
			[]string{
				"29568666068035183841425683795340791879727309630931025356555_4958044192186797981418233587017209679042592862002427381542",
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := func(in, out chan interface{}) {
				for _, v := range tt.data {
					out <- v
				}
			}

			var got []string
			output := func(in, out chan interface{}) {
				for v := range in {
					got = append(got, v.(string))
				}
			}

			jobs := append([]job{input}, tt.jobs...)
			jobs = append(jobs, output)
			ExecutePipeline(jobs...)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("get %v, want %v", got, tt.want)
			}
		})
	}
}
