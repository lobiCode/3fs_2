package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

func main() {

	conn, _ := net.Dial("tcp", "localhost:1234")
	args := os.Args

	s := args[1] + " " + args[2] + "\n"
	conn.Write([]byte(s))
	b := make([]byte, 512)
	conn.Read(b)
	fmt.Println(string(bytes.Split(b, []byte{'\n'})[0]))
	conn.Close()
}
