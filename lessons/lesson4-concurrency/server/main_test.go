package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"testing"
)

func TestHandle(t *testing.T) {
	const str = "test"
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	rw := bufio.NewReadWriter(bufio.NewReader(bytes.NewBufferString(str)), w)
	Handle(rw)
	err := w.Flush()
	if err != nil {
		t.Error(err)
	}
	if b.String() != str+"\n" {
		t.Errorf("expected %s, got %s", str, b.String())
	}
}

func TestServe(t *testing.T) {
	const addr = ":0000"
	const str = "test\n"

	l, err := net.Listen("tcp", addr)
	defer l.Close()
	if err != nil {
		t.Fatal(err)
	}

	serveDone := make(chan struct{})
	lineCh := make(chan string)

	go func() {
		t.Log(Serve(l, func(rw io.ReadWriter) {
			reader := bufio.NewReader(rw)
			line, err := reader.ReadString('\n')
			if err != nil {
				t.Errorf("failed to read from the connection %s", err.Error())
			}
			lineCh <- line
			close(lineCh)
		}))
		close(serveDone)
	}()

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Errorf("failed to dial the connection %s", err.Error())
	}

	// don't leak connections in tests
	defer conn.Close()

	if _, err = conn.Write([]byte(str)); err != nil {
		t.Errorf("failed to write to the connection %s", err.Error())
	}
	line := <-lineCh
	if line != str {
		t.Errorf("expected %s, got %s", str, line)
	}

	l.Close()
	<-serveDone
}
