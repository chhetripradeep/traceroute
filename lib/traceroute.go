package lib

import (
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const MAX_TTL = 64

// Hop contains details about single hop
type Hop struct {
	TTL     int
	Addr    net.Addr
	RTT     time.Duration
	Type    icmp.Type
	Result  string
}

// TraceRoute returns channel of hop info
func TraceRoute(host string) (<-chan Hop, <-chan error) {
	errc := make(chan error, 1)

	destination, err := net.Dial("ip:icmp", host)
	if err != nil {
		log.Println(err)
		errc <- fmt.Errorf("name %s is invalid", host)
		defer close(errc)
		return nil, errc
	}

	ttl := 1
	timeout := time.Second

	out := make(chan Hop)
	go func() {
		defer close(out)
		defer close(errc)

		for {
			hop, err := sendICMPEcho(destination.RemoteAddr(), ttl, ttl, timeout)
			if err != nil {
				errc <- err
				break
			}
			out <- hop
			ttl++
			if hop.Result == "Success" {
				if hop.Type == ipv4.ICMPTypeEchoReply {
					break
				}
				timeout = SetTimeout(hop.RTT)
			}
			if ttl > MAX_TTL {
				return
			}
		}
	}()

	return out, errc
}

// SetTimeout returns the sane timeout value based on rtt
func SetTimeout(t time.Duration) (time.Duration) {
	return t*3 + time.Millisecond*50
}
