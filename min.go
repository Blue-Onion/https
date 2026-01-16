package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)
var ERROR_MALFORMED_REQUEST_LINE =fmt.Errorf("Bad start Line")
var  ERROR_UNSUPPORTED_HTTP_VERSION =fmt.Errorf("Unsupported http req")
const seprator ="\r\n"

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}
type Request struct {
	RequestLine RequestLine
}

func parseRequestLine(s string) (*RequestLine,string ,error){

	i:=strings.Index(s,seprator)
	if i==-1{
		return nil,s,nil
	}

	startofLine:=s[:i]
	restOfMess:=s[i+len(seprator):]
	parts:=strings.Split(startofLine," ")
	if len(parts)!=3{
		return nil,s,ERROR_MALFORMED_REQUEST_LINE
	}

	httpParts:=strings.Split(parts[2],"/")
	fmt.Println(len(httpParts))
	if len(httpParts)!=2||httpParts[0]!="HTTP"||httpParts[1]!="1.1"{
		fmt.Print("this one")
		return nil,restOfMess,ERROR_MALFORMED_REQUEST_LINE
	}
	rl:= &RequestLine{
		Method:parts[0],
		RequestTarget: parts[1],
		HttpVersion: httpParts[1],
	}
	fmt.Println("Fine til here")
	return rl,restOfMess,nil
}
func RequestFromReader(reader io.Reader) (*Request, error) {
	data,err:=io.ReadAll(reader)
	if err!=nil{
		return nil,errors.Join(fmt.Errorf("unable to read all io.reader:"),err)
	}
	str:=string(data)
	rl,_,err:=parseRequestLine(str)
	if err!=nil{
		return nil,errors.Join(fmt.Errorf("unable to read all io.reader:"),err)
	}
	return &Request{
		RequestLine: *rl,
	},nil

}
