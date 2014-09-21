package branchio

import (
	"fmt"
	"log"
	"runtime"
	"testing"
	"time"
)

type printwriter struct{}

func (*printwriter) Write(p []byte) (n int, err error) {
	fmt.Println(string(p))
	return len(p), nil
}

func TestNewWriterSize(t *testing.T) {

	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	a := NewBranchChannelWriter(&printwriter{})
	b := NewBranchChannelWriter(&printwriter{})
	writer := NewWriterSize([]*BranchChannelWriter{a, b}, 10)
	writer.Write([]byte("1111111111111"))
	writer.Write([]byte("1111111111111"))
	writer.Write([]byte("1111111111111"))
	writer.Write([]byte("1111111111111"))
	writer.Write([]byte("1111111111111"))
	writer.Write([]byte("1111111111111"))
	writer.Write([]byte("1111111111111"))
	for writer.lenOutputBuffers() > 0 {
		log.Printf("len: %d\n", writer.lenOutputBuffers())
		time.Sleep(10 * time.Microsecond)
	}
}
