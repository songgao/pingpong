package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type renderer struct {
	printTime bool
	textWidth int
	numIn     int
	once      *sync.Once
	stream    chan item
	in        []chan []byte

	In []chan<- []byte
}

type item struct {
	Index    int
	Content  []byte
	IsStatus bool
}

func NewRender(numIn int, printTime bool) *renderer {
	r := new(renderer)
	r.printTime = printTime
	r.textWidth = getTermWidth() - 9 - numIn - 1
	r.numIn = numIn
	r.stream = make(chan item, 16)
	r.once = new(sync.Once)
	r.in = make([]chan []byte, numIn)
	r.In = make([]chan<- []byte, numIn)
	for i := 0; i < numIn; i++ {
		r.in[i] = make(chan []byte)
		r.In[i] = r.in[i]
	}
	return r
}

func (r *renderer) PrintLegend(titles []string) {
	for i := 0; i < r.numIn; i++ {
		fmt.Fprintf(os.Stdout, rainbowBG("        ", i))
		fmt.Fprintf(os.Stdout, " %s\n", titles[i])
	}
}

func (r *renderer) collector(indexIn int, stream chan item, wg *sync.WaitGroup) {
	ok := true
	var buf []byte
	for ok {
		buf, ok = <-r.in[indexIn]
		if ok && len(buf) > 0 {
			stream <- item{Index: indexIn, Content: buf, IsStatus: false}
		}
	}
	wg.Done()
}

func (r *renderer) Run() {
	r.once.Do(func() {
		wg := new(sync.WaitGroup)
		wg.Add(r.numIn)
		for i := 0; i < r.numIn; i++ {
			go r.collector(i, r.stream, wg)
		}
		go func() {
			var it item
			ok := true
			for ok {
				it, ok = <-r.stream
				if ok {
					r.print(it)
				}
			}
		}()
		wg.Wait()
		close(r.stream)
	})
}

func (r *renderer) PrintStatus(str []byte, indexIn int) {
	r.stream <- item{Index: indexIn, Content: str, IsStatus: true}
}

func flag_init(numIn int) int {
	flag := 0
	for i := uint(0); i < uint(numIn); i++ {
		flag = flag | (1 << i)
	}
	return flag
}

func (r *renderer) print(it item) {
	if it.Content[len(it.Content)-1] == '\n' {
		it.Content = it.Content[:len(it.Content)-1]
	}
	lines := slice(it.Content, r.textWidth)
	starting := true
	for _, line := range lines {
		r.printHeader(os.Stdout, it.Index, starting)
		starting = false
		if it.IsStatus {
			writeBold(os.Stdout, line)
		} else {
			os.Stdout.Write(line)
		}
		fmt.Fprintln(os.Stdout)
	}
}

func (r *renderer) printHeader(w io.Writer, indexIn int, starting bool) {
	if r.printTime {
		if starting {
			t := time.Now()
			fmt.Fprintf(w, "%02d:%02d:%02d ", t.Hour(), t.Minute(), t.Second())
		} else {
			fmt.Fprint(w, "         ")
		}
	}
	for i := 0; i < r.numIn; i++ {
		if indexIn == i {
			fmt.Fprintf(w, rainbowBG(" ", i))
		} else {
			fmt.Fprintf(w, rainbowFG("|", i))
		}
	}
	fmt.Fprint(w, " ")
}

func slice(str []byte, width int) [][]byte {
	var ret_len int
	if 0 == len(str)%width {
		ret_len = len(str) / width
	} else {
		ret_len = len(str)/width + 1
	}
	ret := make([][]byte, ret_len)
	for i := 0; i < ret_len-1; i++ {
		ret[i] = str[width*i : width*(i+1)]
	}
	ret[ret_len-1] = make([]byte, width)
	copy(ret[ret_len-1], str[width*(ret_len-1):])
	return ret
}
