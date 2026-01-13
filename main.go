package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open("message.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	str:=""
	for {
		data:=make([]byte, 8)
		n,err:=f.Read(data)
		if err!=nil{
			log.Fatal(err)
			break
		}
		data=data[:n]
		if i:=bytes.IndexByte(data,'\n');i!=-1{
			str+=string(data[:i]);

			data=data[i+1:]
			fmt.Println(str)
			str=""
		}
		str+=string(data)
	}
	if len(str)!=0{
		fmt.Println(str)
	}
}
