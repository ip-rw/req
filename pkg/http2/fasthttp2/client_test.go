package fasthttp2

import (
	"bufio"
	"bytes"
	"io"
	"math/rand"
	"testing"

	"github.com/ip-rw/req/pkg/http2"
)

func TestClientWriteOrder(t *testing.T) {
	bf := bytes.NewBuffer(nil)

	c := &http2.Client{}
	c.writer = make(chan *Frame, 1)
	c.bw = bufio.NewWriter(bf)

	go c.writeLoop()

	framesToTest := 32

	id := uint32(1)
	frames := make([]*Frame, 0, framesToTest)

	for i := 0; i < framesToTest; i++ {
		fr := AcquireFrame()
		fr.SetStream(id)
		id += 2
		frames = append(frames, fr)
	}

	for len(frames) > 0 {
		i := rand.Intn(len(frames))

		c.writeFrame(frames[i])
		frames = append(frames[:i], frames[i+1:]...)
	}

	br := bufio.NewReader(bf)
	fr := AcquireFrame()

	expected := uint32(1)
	for i := 0; i < framesToTest; i++ {
		_, err := fr.ReadFrom(br)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatalf("error reading frame: %s", err)
		}

		if fr.Stream() != expected {
			t.Fatalf("Unexpected id: %d <> %d", fr.Stream(), expected)
		}

		expected += 2
	}

	close(c.writer)
}
