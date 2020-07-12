package main

import (
	"fmt"
	"log"
	"os"
	"traceroute/lib"
)

func main() {
	l := log.New(os.Stdout, "", log.Llongfile|log.Ldate)
	fmt.Println("Starting")
	c, e := lib.TraceRoute("yahoo.com")
	select {
	case r := <-c:
		l.Println(r)
	case err := <- e:
		l.Println(err)
	}
}