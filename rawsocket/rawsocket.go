package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"./myethernet"

	"github.com/mdlayher/raw"
)

func main() {
	go RecTest()
	time.Sleep(time.Duration(time.Second * 2))
	ifi, err := net.InterfaceByName("wlp61s0")
	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
	}
	// Open a raw socket using same EtherType as our frame.
	c, err := raw.ListenPacket(ifi, 0x7676, nil)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer c.Close()
	// Marshal a frame to its binary format.
	fmt.Printf("%v\n", ifi.HardwareAddr)
	f := myethernet.Frame{
		Source:      ifi.HardwareAddr,
		Destination: net.HardwareAddr{0xa0, 0x88, 0x69, 0x16, 0xda, 0xb4},
		EtherType:   0x7676,
		Payload:     []byte("hello world")}
	b, err := f.MarshalBinary()
	if err != nil {
		log.Fatalf("failed to marshal frame: %v", err)
	}
	// Broadcast the frame to all devices on our network segment.
	addr := &raw.Addr{HardwareAddr: net.HardwareAddr{0xa0, 0x88, 0x69, 0x16, 0xda, 0xb4}}
	if _, err := c.WriteTo(b, addr); err != nil {
		log.Fatalf("failed to write frame: %v", err)
	}
	for {
	}
}

func RecTest() {
	fmt.Println("Start to rec")
	ifi, err := net.InterfaceByName("wlp61s0")
	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
	}
	// Open a raw socket using same EtherType as our frame.
	c, err := raw.ListenPacket(ifi, 0x7676, nil)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer c.Close()
	// Accept frames up to interface's MTU in size.
	b := make([]byte, ifi.MTU)
	var f myethernet.Frame
	// Keep reading frames.
	for {
		n, addr, err := c.ReadFrom(b)
		if err != nil {
			log.Fatalf("failed to receive message: %v", err)
		}
		// Unpack Ethernet frame into Go representation.
		if err := (&f).UnmarshalBinary(b[:n]); err != nil {
			log.Fatalf("failed to unmarshal ethernet frame: %v", err)
		}
		// Display source of message and message itself.
		log.Printf("[%s] %s", addr.String(), string(f.Payload))
	}
}
