package io

import (
	"io/ioutil"
	"os"
	"testing"

	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"

	testu "github.com/ipfs/go-ipfs/unixfs/test"
)

func TestBasicRead(t *testing.T) {
	dserv := testu.GetDAGServ()
	inbuf, node := testu.GetRandomNode(t, dserv, 1024)
	ctx, closer := context.WithCancel(context.Background())
	defer closer()

	reader, err := NewDagReader(ctx, node, dserv)
	if err != nil {
		t.Fatal(err)
	}

	outbuf, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	err = testu.ArrComp(inbuf, outbuf)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSeekAndRead(t *testing.T) {
	dserv := testu.GetDAGServ()
	inbuf := make([]byte, 256)
	for i := 0; i <= 255; i++ {
		inbuf[i] = byte(i)
	}

	node := testu.GetNode(t, dserv, inbuf)
	ctx, closer := context.WithCancel(context.Background())
	defer closer()

	reader, err := NewDagReader(ctx, node, dserv)
	if err != nil {
		t.Fatal(err)
	}

	for i := 255; i >= 0; i-- {
		reader.Seek(int64(i), os.SEEK_SET)

		if reader.Offset() != int64(i) {
			t.Fatal("expected offset to be increased by one after read")
		}

		out := readByte(t, reader)

		if int(out) != i {
			t.Fatalf("read %d at index %d, expected %d", out, i, i)
		}

		if reader.Offset() != int64(i+1) {
			t.Fatal("expected offset to be increased by one after read")
		}
	}
}

func TestRelativeSeek(t *testing.T) {
	dserv := testu.GetDAGServ()
	ctx, closer := context.WithCancel(context.Background())
	defer closer()

	inbuf := make([]byte, 1024)

	for i := 0; i < 256; i++ {
		inbuf[i*4] = byte(i)
	}

	inbuf[1023] = 1 // force the reader to be 1024 bytes
	node := testu.GetNode(t, dserv, inbuf)

	reader, err := NewDagReader(ctx, node, dserv)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 256; i++ {
		if reader.Offset() != int64(i*4) {
			t.Fatalf("offset should be %d, was %d", i*4, reader.Offset())
		}
		out := readByte(t, reader)
		if int(out) != i {
			t.Fatalf("expected to read: %d at %d, read %d", i, reader.Offset()-1, out)
		}
		if i != 255 {
			_, err := reader.Seek(3, os.SEEK_CUR)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	_, err = reader.Seek(4, os.SEEK_END)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 256; i++ {
		if reader.Offset() != int64(1020-i*4) {
			t.Fatalf("offset should be %d, was %d", 1020-i*4, reader.Offset())
		}
		out := readByte(t, reader)
		if int(out) != 255-i {
			t.Fatalf("expected to read: %d at %d, read %d", 255-i, reader.Offset()-1, out)
		}
		reader.Seek(-5, os.SEEK_CUR) // seek 4 bytes but we read one byte every time so 5 bytes
	}

}

func readByte(t testing.TB, reader *DagReader) byte {
	out := make([]byte, 1)
	c, err := reader.Read(out)

	if c != 1 {
		t.Fatal("reader should have read just one byte")
	}
	if err != nil {
		t.Fatal(err)
	}

	return out[0]
}