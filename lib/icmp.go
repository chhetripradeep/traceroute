package lib

import (
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"math/rand"
	"net"
	"time"
)

const (
	MTU = 1500
	LISTEN_ADDRESS = "0.0.0.0"
)

// send an icmp echo message
func sendICMPEcho(destination net.Addr, sequence int, ttl int, timeout time.Duration) (hop Hop, err error) {
	// start listening for response
	connection, err := icmp.ListenPacket("ip4:icmp", LISTEN_ADDRESS)
	if err != nil {
		return Hop{}, err
	}

	// create icmp echo packet
	message, err := createICMPEcho(sequence)
	if err != nil {
		return Hop{}, err
	}

	// set ttl for the icmp echo packet
	connection.IPv4PacketConn().SetTTL(ttl)

	// start recording time for measuring round-trip time
	start := time.Now()

	// send icmp echo message over the connection
	_, err = connection.WriteTo(message, destination)
	if err != nil {
		return Hop{}, err
	}

	// allocate slice for response
	response := make([]byte, MTU)
	err = connection.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return Hop{}, err
	}

	// read response from the connection
	_, peer, err := connection.ReadFrom(response)
	if err != nil {
		return Hop{TTL: ttl, Result: "Failure"}, err
	}

	// calculate round-trip time
	rtt := time.Since(start)

	// parse response message
	reply, err := icmp.ParseMessage(1, response)
	if err != nil {
		return Hop{TTL: ttl, Result: "Failure"}, err
	}

	hop = Hop{
		TTL: ttl,
		Addr: peer,
		RTT: rtt,
		Type: reply.Type,
		Result: "Success",
	}

	return
}

// create an icmp echo message
func createICMPEcho(sequence int) (request []byte, err error) {
	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   rand.Int(),
			Seq:  sequence,
			Data: []byte(""),
		},
	}
	return message.Marshal(nil)
}