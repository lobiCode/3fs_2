package qman

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"math"
	"math/big"
	"net"
	"strconv"
)

const (
	R   = 2.236067977499789696409173668731276235440618359611525724270
	PHI = (1 + R) / 2
	Phi = (1 - R) / 2
)

var (
	FIBONACCI  = []byte("Fibonacci")
	REVERSE    = []byte("ReverseText")
	ENCODER    = []byte("TextEncoder")
	ARITHMETIC = []byte("BasicArithmetic")
)

type Resolver interface {
	R()
}

type Job struct {
	Res Resolver
}

type Worker struct {
	JobQueue chan Job
}

func CreateWorker(jobQueue chan Job) Worker {
	return Worker{jobQueue}
}

func (w *Worker) Start() {

	go func() {
		for {
			select {
			case job := <-w.JobQueue:
				job.Res.R()
			}
		}
	}()
}

type Fibonacci struct {
	B    []byte
	Conn net.Conn
}

func (f Fibonacci) R() {

	n, err := strconv.Atoi(string(f.B))
	if err != nil {
		// TODO err msq, close, return
	}
	// TODO handlaj ce je n prevelik

	nfloat := float64(n)
	fib := (math.Pow(PHI, nfloat) - math.Pow(Phi, nfloat)) / R
	s := strconv.FormatFloat(fib, 'f', 0, 64)
	s = s + "\n"
	f.Conn.Write([]byte(s))
	f.Conn.Close()
}

type ReverseText struct {
	B    []byte
	Conn net.Conn
}

func (rt ReverseText) R() {

	for i := 0; i < len(rt.B)/2; i++ {
		j := len(rt.B) - 1 - i
		rt.B[i], rt.B[j] = rt.B[j], rt.B[i]
	}

	rt.B = append(rt.B, '\n')
	rt.Conn.Write(rt.B)
	rt.Conn.Close()
}

type TextEncoder struct {
	B    []byte
	Conn net.Conn
}

func (te TextEncoder) R() {

	b, err := bcrypt.GenerateFromPassword(te.B, bcrypt.DefaultCost)
	if err != nil {
		// TODO
	}

	te.Conn.Write(b)
	te.Conn.Close()

}

type BasicArithmetic struct {
	B    []byte
	Conn net.Conn
}

func (ba BasicArithmetic) R() {

	var i int
	var b byte
	for i, b = range ba.B {
		if findOpera(b) {
			break
		}
	}

	if i == len(ba.B) {
		writeAndClose([]byte("wrong arg"), ba.Conn)
		return
	}

	s1 := string(ba.B[:i])
	s2 := string(ba.B[i+1:])

	i1, err := strconv.Atoi(s1)
	if err != nil {
		writeAndClose([]byte(err.Error()), ba.Conn)
		return
	}
	i2, err := strconv.Atoi(s2)
	if err != nil {
		writeAndClose([]byte(err.Error()), ba.Conn)
		return
	}

	if b == '/' {
		fb, err := div(i1, i2)
		if err == nil {
			writeAndClose([]byte(fb.String()), ba.Conn)
			return
		} else {
			writeAndClose([]byte(err.Error()), ba.Conn)
			return
		}
	}

	ib1 := big.NewInt(int64(i1))
	ib2 := big.NewInt(int64(i2))
	switch {
	case b == '+':
		ib1.Add(ib1, ib2)
	case b == '-':
		ib1.Sub(ib1, ib2)
	case b == '*':
		ib1.Mul(ib1, ib2)
	}

	writeAndClose([]byte(ib1.String()), ba.Conn)

}

func findOpera(b byte) bool {

	if b == '+' || b == '-' || b == '*' || b == '/' {
		return true
	}

	return false
}

func div(i, j int) (bf *big.Float, err error) {

	if j == 0 {
		return bf, errors.New("can't divide by zero")
	}

	if i == 0 {
		return big.NewFloat(0), nil
	}

	fb1 := big.NewFloat(float64(i))
	fb2 := big.NewFloat(float64(j))

	fb1.Quo(fb1, fb2)

	return fb1, nil
}

func writeAndClose(b []byte, conn net.Conn) {

	conn.Write(b)
	conn.Close()
}
