package agent

import (
	"bytes"
	"io"
	"sync"

	"github.com/armon/circbuf"
)

type Buffer struct {
	buf    *circbuf.Buffer
	pipes  map[*PipeReader]*io.PipeWriter
	muBuf  sync.Mutex
	muPipe sync.RWMutex
}

func NewBuffer(size int64) (*Buffer, error) {
	cbuf, err := circbuf.NewBuffer(size)
	if err != nil {
		return nil, err
	}
	return &Buffer{
		buf:   cbuf,
		pipes: make(map[*PipeReader]*io.PipeWriter),
	}, nil
}

func (b *Buffer) CircularBytes() []byte {
	b.muBuf.Lock()
	by := b.buf.Bytes()
	b.muBuf.Unlock()
	return by
}

func (b *Buffer) Status() (int, int64) {
	if b == nil {
		return 0, 0
	}

	b.muBuf.Lock()
	tw := b.buf.TotalWritten()
	n := len(b.pipes)
	b.muBuf.Unlock()
	return n, tw
}

func (b *Buffer) Write(by []byte) (int, error) {
	b.muBuf.Lock()
	n, err := b.buf.Write(by)
	b.muBuf.Unlock()
	if err != nil {
		return n, err
	}

	b.muPipe.RLock()
	defer b.muPipe.RUnlock()
	for _, pw := range b.pipes {
		pn, err := pw.Write(by)
		if err != nil {
			return n, err
		}
		if pn != n {
			return n, io.ErrShortWrite
		}
	}

	return n, err
}

func (b *Buffer) Release(r *PipeReader) {
	b.muPipe.Lock()
	delete(b.pipes, r)
	b.muPipe.Unlock()
}

func (b *Buffer) Pipe() io.ReadCloser {
	pr, pw := io.Pipe()
	pr2 := &PipeReader{
		Reader: io.MultiReader(bytes.NewBuffer(b.CircularBytes()), pr),
		pr:     pr,
		buf:    b,
	}

	b.muPipe.Lock()
	b.pipes[pr2] = pw
	b.muPipe.Unlock()
	return pr2
}

func (b *Buffer) Close() error {
	b.muPipe.Lock()
	for pr := range b.pipes {
		pr.pr.Close()
		delete(b.pipes, pr)
	}
	b.muPipe.Unlock()
	return nil
}

type PipeReader struct {
	io.Reader                // multi reader containing circbuf + pipe reader
	pr        *io.PipeReader // keep ref so it can be closed directly
	buf       *Buffer
}

func (r *PipeReader) Close() error {
	err := r.pr.Close()
	r.buf.Release(r)
	return err
}
