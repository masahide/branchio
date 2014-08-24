package branchio

const (
	defaultChannelSize = 100
)

type Writer struct {
	err           error
	outputBuffers [](chan<- []byte)
}

type BranchWriter func(<-chan []byte)

func NewWriterSize(branchWriters []BranchWriter, size int) *Writer {
	outputBuffers := [](chan<- []byte){}
	if size <= 0 {
		size = defaultChannelSize
	}
	// Is it already a Writer?
	for _, branchWriter := range branchWriters {
		outputBuffer := make(chan []byte, size)
		outputBuffers = append(outputBuffers, outputBuffer)
		go branchWriter(outputBuffer)
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

func NewWriter(branchWriters []BranchWriter) *Writer {
	return NewWriterSize(branchWriters, defaultChannelSize)
}
