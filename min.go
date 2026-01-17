package request

import (
	"bytes"
	"fmt"
	"io"
)

type RequestLine struct {
	HttpVersion   string
	Method        string
	RequestTarget string
}
type Request struct {
	RequestLine RequestLine
	state       parseState
}
type parseState string

const (
	StateInit  parseState = "Init"
	StateErorr parseState = "error"
	StateDone  parseState = "Done"
)

var seprator = []byte("\r\n")

var ErrorMalformedRequestLine = fmt.Errorf("Bad start Line")
var ErrorUnsupportedHttpVersion = fmt.Errorf("Unsupported http req")
var ErrorRequestInErorrState = fmt.Errorf("request in error state")

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}
func (req *Request) parse(data []byte) (int, error) {

	read := 0
outer:
	for {

		switch req.state {
		case StateErorr:
			return 0,ErrorRequestInErorrState
		case StateInit:
			rl, n, err := parseRequestLine(data)
			if err != nil {
				req.state=StateErorr
				return 0, err
			}
			if n == 0 {
				break outer
			}
			req.RequestLine = *rl
			read += n
			req.state = StateDone
		case StateDone:
			break outer
		}
	}
	return read, nil

}
func (req *Request) done() bool {
	return req.state == StateErorr||req.state==StateDone
}

func parseRequestLine(s []byte) (*RequestLine, int, error) {

	i := bytes.Index(s, seprator)
	if i == -1 {

		return nil, 0, nil
	}
	startOfLine := s[:i]
	read := i + len(seprator)
	parts := bytes.Split(startOfLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrorMalformedRequestLine
	}
	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ErrorUnsupportedHttpVersion
	}
	return &RequestLine{
		Method:        string(parts[0]),
		HttpVersion:   string(httpParts[1]),
		RequestTarget: string(parts[1]),
	}, read, nil
}
func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err

		}
		fmt.Println(n)
		bufLen += n
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}
	return request, nil
}
