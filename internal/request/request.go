package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)
type RequestLine struct{
	HttpVersion string
	Method string
	RequestTarget string
}
type Request struct{
	RequestLine RequestLine
}
const seprator="\r\n"
var ERROR_MALFORMED_REQUEST_LINE =fmt.Errorf("Bad start Line")
var  ERROR_UNSUPPORTED_HTTP_VERSION =fmt.Errorf("Unsupported http req")
func parseRequestLine(s string) (*RequestLine,string,error){
	i:=strings.Index(s,seprator)
	if i==-1{
		return nil,s,ERROR_MALFORMED_REQUEST_LINE
	}
	startOfLine:=s[:i]
	endLine:=s[i+len(seprator):]
	parts:=strings.Split(startOfLine," ")
	if len(parts)!=3{
		return nil,startOfLine,ERROR_MALFORMED_REQUEST_LINE
	}
	httpParts:=strings.Split(parts[2],"/")
	if len(httpParts)!=2||httpParts[0]!="HTTP"||httpParts[1]!="1.1"{
		return nil,startOfLine,ERROR_MALFORMED_REQUEST_LINE
	}
	return &RequestLine{
		Method: parts[0],
		HttpVersion: httpParts[1],
		RequestTarget: parts[1],

	},endLine,nil
}
func RequestFromReader(reader io.Reader) (*Request, error){
	data,err:=io.ReadAll(reader);
	if err!=nil{
		return nil,errors.Join(fmt.Errorf("unable to read all io.reader:"),err)
	}
	str:=string(data)
	rl,_,err:=parseRequestLine(str)
	if err!=nil{
		return nil,errors.Join(fmt.Errorf("unable to read all io.reader:"),err)
	}

	
	return &Request{RequestLine:*rl},nil
}