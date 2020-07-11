package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/chhetripradeep/traceroute"
)

func main() {
	host := flag.Arg(0)
	hops, errs := TraceRoute(host)
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