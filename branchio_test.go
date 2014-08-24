package branchio

import (
	"log"
	"runtime"
	"testing"
	"time"
)

func hoge(buf <-chan []byte) {
	for {
		log.Printf("hoge:%s\n", <-buf)
	}
}
func fuga(buf <-chan []byte) {
	for {
		log.Printf("fuga:%s\n", <-buf)
	}
}

func TestNewWriterSize(t *testing.T) {

	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	writer := NewWriterSize([]BranchWriter{hoge, fuga}, 10)
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
