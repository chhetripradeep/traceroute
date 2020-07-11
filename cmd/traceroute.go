package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/chhetripradeep/traceroute"
)

func usage()  {
	fmt.Println("You must specify one host as argument.")
	flag.Usage()
}

func main() {
	hosts := flag.Args()
	if hosts == nil || len(hosts) == 0 {
		usage()
		return
	}
	hops, errs := TraceRoute(hosts[0])
	for {
		select {
		case err, ok := <-errs:
			if !ok {
				return
			}
			printErr(err)
		case hop, ok := <-hops:
			if !ok {
				return
			}
			printHop(hop)
			fmt.Println()
		}
	}
}

func printHop(hop Hop) {
	fmt.Printf("%2d ", hop.Number)
	if hop.Result == "Success" {
		fmt.Println(hop.Addr.String())
	} else {
		fmt.Println("*")
	}
}

func printErr(err error) {
	if strings.Contains(err.Error(), "Operation not permitted") {
		fmt.Println("Try running the program as root, or give it permission to use the network like this:")
		fmt.Println("$ sudo setcap cap_net_raw+ep traceroute")
	} else {
		fmt.Println(err)
	}
}