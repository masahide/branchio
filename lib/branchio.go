package lib

import "io"

const (
	defaultChannelSize = 100
)

type BranchChannelWriter struct {
	writer io.Writer
	errc   <-chan error
	in     chan []byte
}

func (w *BranchChannelWriter) channelToWriter(in <-chan []byte) <-chan error {

	errc := make(chan error, 1)
	go func() {
		for p := range in {
			n, err := w.writer.Write(p)
			if err != nil {
				errc <- err
				return
			}
			if n != len(p) {
				errc <- io.ErrShortWrite
				return
			}
		}
		errc <- nil
	}()
	return errc
}

type BranchWriter struct {
	branchChannelWriters []*BranchChannelWriter
}

func NewBranchWriter(bufSize int, writers ...io.Writer) *BranchWriter {
	branchChannelWriters := make([]*BranchChannelWriter, 0, len(writers))
	if bufSize <= 0 {
		bufSize = defaultChannelSize
	}
	for _, w := range writers {
		bcw := BranchChannelWriter{
			writer: w,
			in:     make(chan []byte, bufSize),
		}
		bcw.errc = bcw.channelToWriter(bcw.in)
		branchChannelWriters = append(branchChannelWriters, &bcw)
	}
	return &BranchWriter{branchChannelWriters}
}

func (t *BranchWriter) Write(p []byte) (n int, err error) {
	for _, bcw := range t.branchChannelWriters {
		buf := make([]byte, len(p))
		copy(buf, p)
		bcw.in <- buf
	}
	return len(p), nil
}
func (t *BranchWriter) Close() {
	for _, bcw := range t.branchChannelWriters {
		close(bcw.in)
	}
}

func (t *BranchWriter) Merge() []error {
	errs := make([]error, 0, len(t.branchChannelWriters))
	for _, bcw := range t.branchChannelWriters {
		errs = append(errs, <-bcw.errc)
	}
	return errs
}
