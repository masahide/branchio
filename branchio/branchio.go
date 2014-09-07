package branchio

import (
	"io"
	"log"
)

const (
	defaultChannelSize = 100
)

type Writer struct {
	err           error
	outputBuffers map[int](chan<- []byte)
}

type BranchChannelWriter struct {
	writer io.Writer
}

func NewBranchChannelWriter(writer io.Writer) *BranchChannelWriter {
	return &BranchChannelWriter{
		writer: writer,
	}

}

func (this *BranchChannelWriter) Writer(c <-chan []byte) {
	for b := range c {
		for n := 0; len(b) > n; {
			length, err := this.writer.Write(b[n:])
			if err != nil {
				log.Printf("Error Branch Write: %v", err)
				break
			}
			n += length
		}
	}
}

func NewWriterSize(branchWriters []*BranchChannelWriter, size int) *Writer {
	outputBuffers := map[int](chan<- []byte){}
	if size <= 0 {
		size = defaultChannelSize
	}
	// Is it already a Writer?
	for _, branchWriter := range branchWriters {
		outputBuffer := make(chan []byte, size)
		outputBuffers[len(outputBuffers)] = outputBuffer
		go branchWriter.Writer(outputBuffer)
	}
	return &Writer{
		outputBuffers: outputBuffers,
	}
}

func (b *Writer) Write(p []byte) (nn int, err error) {
	for _, output := range b.outputBuffers {
		output <- p
	}
	nn = len(p)
	return
}

func (b *Writer) lenOutputBuffers() int {
	length := 0
	for _, output := range b.outputBuffers {
		length += len(output)
	}
	return length
}

func NewWriter(branchWriters []*BranchChannelWriter) *Writer {
	return NewWriterSize(branchWriters, defaultChannelSize)
}
