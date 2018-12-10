package rawsocket

import (
	"fmt"
	"log"
	"net"
	"sync"

	"./myethernet"

	"github.com/mdlayher/raw"
)

// func main() {

// 	ifaceName := "wlp4s0"
// 	dstMac := net.HardwareAddr{0x34, 0xe6, 0xad, 0x09, 0xc6, 0x3f}
// 	myPayload := []byte("Hello World")

// 	var llwg sync.WaitGroup
// 	for i := 0; i < 10; i++ {
// 		llwg.Add(1)
// 		go SendLinkLayer(ifaceName, dstMac, 0x7676, myPayload, &llwg)
// 	}
// 	RecLinkLayer(ifaceName, 0x7676)
// 	llwg.Wait()
// }

//SendLinkLayer send a packet via link layer socket
func SendLinkLayer(ifaceName string, dstMac net.HardwareAddr, ethType myethernet.EtherType, myPayload []byte, llwg *sync.WaitGroup) {
	ifi, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
	}
	// Open a raw socket using same EtherType as our frame.
	c, err := raw.ListenPacket(ifi, uint16(ethType), nil)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer c.Close()
	defer llwg.Done()
	// Marshal a frame to its binary format.
	fmt.Printf("%v\n", ifi.HardwareAddr)
	f := myethernet.Frame{
		Source:      ifi.HardwareAddr,
		Destination: dstMac,
		EtherType:   ethType, //0x7676,
		Payload:     myPayload,
	}
	b, err := f.MarshalBinary()
	if err != nil {
		log.Fatalf("failed to marshal frame: %v", err)
	}
	// Broadcast the frame to all devices on our network segment.
	addr := &raw.Addr{HardwareAddr: dstMac}
	if _, err := c.WriteTo(b, addr); err != nil {
		log.Fatalf("failed to write frame: %v", err)
	}
}

//RecLinkLayer link layer receiver
func RecLinkLayer(ifaceName string, ethType myethernet.EtherType) {
	fmt.Println("Start to rec")
	ifi, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
	}
	// Open a raw socket using same EtherType as our frame.
	c, err := raw.ListenPacket(ifi, uint16(ethType), nil)
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
		go RecHandler(f, addr, b, n)
	}
}

//RecHandler handler
func RecHandler(f myethernet.Frame, addr net.Addr, b []byte, n int) {
	// Unpack Ethernet frame into Go representation.
	if err := (&f).UnmarshalBinary(b[:n]); err != nil {
		log.Fatalf("failed to unmarshal ethernet frame: %v", err)
	}
	// Display source of message and message itself.
	log.Printf("[%s] %s", addr.String(), string(f.Payload))
}
