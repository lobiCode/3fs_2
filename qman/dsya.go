package qman

import (
	"errors"
	"go/parser"
	"math/big"
)

// TODO daj raje vse v string
type Stack struct {
	S []interface{}
}

func (s *Stack) Push(i interface{}) {
	s.S = append(s.S, i)
}

func (s *Stack) Pop() (interface{}, bool) {
	i := len(s.S)
	if i == 0 {
		return nil, false
	}

	e := s.S[i-1]
	s.S = append(s.S[:0], s.S[:i-1]...)
	return e, true

}

func (s *Stack) PopFirst() (interface{}, bool) {
	i := len(s.S)
	if i == 0 {
		return nil, false
	}

	e := s.S[0]
	s.S = append(s.S[:0], s.S[1:]...)
	return e, true
}

func Dsya(b []byte) (string, error) {

	// TODO remove
	_, err := parser.ParseExpr(string(b))
	if err != nil {
		return "", errors.New("wrong arg")
	}
	for _, e := range b {
		if e > 57 || e < 40 || e == 44 || e == 46 {
			return "", errors.New("wrong arg")
		}
	}
	//

	operator := Stack{S: make([]interface{}, 0)}
	output := Stack{S: make([]interface{}, 0)}

	for i := 0; i < len(b); i++ {
		e := b[i]
		le := len(operator.S)

		switch {
		case e < 58 && e > 47:
			i = i + pushFloat(&output, b[i:])
		case e == '(':
			operator.Push(e)
		case e == ')':
			discRightPar(&operator, &output)
		case le == 0 || operator.S[le-1].(byte) == '(' || e == '(':
			operator.Push(e)
		case (e == '*' || e == '/') && (operator.S[le-1].(byte) == '+' || operator.S[le-1].(byte) == '-'):
			operator.Push(e)
		default:
			popaj(&operator, &output, e)
		}
	}

	for e, ok := operator.Pop(); ok; e, ok = operator.Pop() {
		output.Push(e)
	}

	for e, ok := output.PopFirst(); ok; e, ok = output.PopFirst() {
		switch t := e.(type) {
		case byte:
			cal(&operator, t)
		default:
			operator.Push(t)
		}
	}

	t, _ := operator.Pop()
	fb, _ := t.(*big.Float)

	return fb.String(), nil
}

func cal(s *Stack, e byte) {

	t1, _ := s.Pop()
	t2, _ := s.Pop()
	fb1, _ := t1.(*big.Float)
	fb2, _ := t2.(*big.Float)

	switch e {
	case '-':
		fb2.Sub(fb2, fb1)
	case '+':
		fb2.Add(fb1, fb2)
	case '/':
		//TODO 0
		fb2.Quo(fb2, fb1)
	default:
		//TODO 0
		fb2.Mul(fb2, fb1)
	}
	s.Push(fb2)
}

func pushFloat(s *Stack, b []byte) int {

	var i int
	for i = 0; i < len(b); i++ {
		e := b[i]
		if e > 57 || e < 48 {
			break
		}
	}

	bi := new(big.Int)
	bi.SetString(string(b[:i]), 10)
	bf := new(big.Float)
	bf.SetInt(bi)
	s.Push(bf)
	return i - 1
}

func discRightPar(operator, output *Stack) {

	for {
		t, ok := operator.Pop()
		e := t.(byte)
		if e == '(' || !ok {
			break
		}
		output.Push(e)
	}
}

// :D
func popaj(operator, output *Stack, b byte) {

	if b == '*' || b == '/' {
		e, _ := operator.Pop()
		output.Push(e)
		operator.Push(b)
		return
	}

	for {
		e, ok := operator.Pop()
		if !ok {
			operator.Push(b)
			return
		} else if e.(byte) == '(' {
			operator.Push(e)
			operator.Push(b)
			return
		} else {
			output.Push(e)
		}
	}
}
