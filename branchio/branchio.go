package branchio

import (
	"io"
	"log"
)

const (
	defaultChannelSize = 100
)

type BranchChannelWriter struct {
	writer io.Writer
}

func NewBranchChannelWriter(writer io.Writer) *BranchChannelWriter {
	return &BranchChannelWriter{
		writer: writer,
	}

}

func (this *BranchChannelWriter) ChannelToWriter(c <-chan []byte, workCounter chan bool) {

	for b := range c {
		log.Printf("ChannelToWriter workCounter cap:%v", cap(workCounter))
		workCounter <- true
		for n := 0; len(b) > n; {
			length, err := this.writer.Write(b[n:])
			if err != nil {
				log.Printf("Error Branch Write: %v", err)
				break
			}
			n += length
		}
		<-workCounter
	}
}

func NewWriterSize(branchWriters []*BranchChannelWriter, size int) *Writer {
	outputBuffers := map[int](chan<- []byte){}
	if size <= 0 {
		size = defaultChannelSize
	}
	num := len(branchWriters)
	workCounter := make(chan bool, num)
	log.Printf("workCounter cap:%v", cap(workCounter))

	// Is it already a Writer?
	for _, branchWriter := range branchWriters {
		outputBuffer := make(chan []byte, size)
		outputBuffers[len(outputBuffers)] = outputBuffer
		go branchWriter.ChannelToWriter(outputBuffer, workCounter)
	}
	return &Writer{
		outputBuffers: outputBuffers,
		workCounter:   workCounter,
	}
}

type Writer struct {
	err           error
	outputBuffers map[int](chan<- []byte)
	workCounter   chan bool
}

func (b *Writer) CountWorkers() int {
	count := 0
	for _, output := range b.outputBuffers {
		count += len(output)
	}
	count += len(b.workCounter)
	return count

}

func (b *Writer) Write(p []byte) (nn int, err error) {
	for _, output := range b.outputBuffers {
		output <- p
		log.Printf("Writer.Write output len:%v, workCounter cap:%v", len(output), cap(b.workCounter))
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
