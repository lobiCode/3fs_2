package main

import (
	"bytes"
	"net"
	_ "strconv"
	_ "strings"
)

const (
	N      = 4
	BUFFER = 100
)

var (
	JobQueue  chan Job
	FIBONACCI = []byte("Fibonacci")
	REVERSE   = []byte("ReverseText")
	ENCODER   = []byte("TextEncoder")
)

func main() {

	ln, err := net.Listen("tcp", ":1234")

	if err != nil {
		// TODO
	}

	Init()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue

		}
		dispatchers(conn)
	}

}

func Init() {

	JobQueue = make(chan Job, BUFFER)

	for i := 0; i < N; i++ {
		w := CreateWorker(JobQueue)
		w.Start()
	}
}

func dispatchers(conn net.Conn) {

	var j Job
	b := make([]byte, 512)
	conn.Read(b) // TODO len n,

	req := bytes.Split(b, []byte{'\n'})
	req = bytes.Split(req[0], []byte{' '})
	// TODO len(req) < 2  close, return

	switch {
	case bytes.Compare(req[0], FIBONACCI) == 0:
		j = Job{conn, Fibonacci{req[1]}}
	case bytes.Compare(req[0], REVERSE) == 0:
		j = Job{conn, ReverseText{req[1]}}
	case bytes.Compare(req[0], ENCODER) == 0:
		j = Job{conn, TextEncoder{req[1]}}
	default:
		b = []byte("is unkowon resolver")
		req = append(req, b)
		conn.Write(bytes.Join(req, []byte{' '}))
		conn.Close()
		return
	}
	JobQueue <- j
}
