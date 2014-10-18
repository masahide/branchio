package lib

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"testing"
)

type printwriter struct{}

func (*printwriter) Write(p []byte) (n int, err error) {
	fmt.Println(string(p))
	return len(p), nil
}

func TestNewWriterSize(t *testing.T) {

	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	a, _ := os.Create("a.txt")
	b, _ := os.Create("b.txt")
	w := NewBranchWriter(0, a, b)
	for i := 0; i < 10; i++ {
		f, err := os.Open("in.txt")
		if err != nil {
			log.Println(err)
		}
		io.Copy(w, f)
		f.Close()
	}
	w.Close()
	log.Printf("errs:%v", w.Merge())
}
