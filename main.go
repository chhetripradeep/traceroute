package main

import (
	"fmt"
	"log"
	"os"
	"traceroute/lib"
	"flag"
)

func main() {
	host := flag.String("host", "127.0.0.1", "hostname/ip to traceroute to")
	l := log.New(os.Stdout, "", log.Llongfile|log.Ldate)
	fmt.Println("Starting")
	c, e := lib.TraceRoute(*host)
	select {
	case r := <-c:
		l.Println(r)
	case err := <- e:
		l.Println(err)
	}
}